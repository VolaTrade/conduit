package service

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/google/wire"
	"github.com/volatrade/candles/internal/cache"
	"github.com/volatrade/candles/internal/client"
	"github.com/volatrade/candles/internal/models"
	"github.com/volatrade/candles/internal/stats"

	"github.com/volatrade/candles/internal/storage"
)

const rootWsURI string = "stream.binance.com:9443"

var (
	Module = wire.NewSet(
		New,
	)
)

type (
	Service interface {
		Init() error
		ConcurrentTickDataCollection()
		CheckForDatabasePriveleges()
	}

	TickersService struct {
		cache     *cache.TickersCache
		store     *storage.ConnectionArray
		exch      *client.ApiClient
		statsd    *stats.StatsD
		writeToDB bool
	}
)

func New(arr *storage.ConnectionArray, ch *cache.TickersCache, cl *client.ApiClient, stats *stats.StatsD) *TickersService {
	return &TickersService{cache: ch, store: arr, exch: cl, statsd: stats, writeToDB: false}
}

func (ts *TickersService) CheckForDatabasePriveleges() {

	for {

		if _, err := os.Stat("start"); err == nil {
			log.Println("making connections to DB NOW")
			ts.store.MakeConnections()
			ts.writeToDB = true
			log.Println("Purging cache")
			err := ts.store.Arr[0].PurgeCache(ts.cache)
			if err != nil {
				panic(err)
			}
			return
		}

	}
}

func (ts *TickersService) Init() error {

	tradingCryptosList, err := ts.exch.GetActiveBinanceExchangePairs()
	if err != nil {
		return err
	}

	for _, val := range tradingCryptosList {
		temp := val.(map[string]interface{}) //type casting
		id := strings.ToLower(temp["symbol"].(string))

		if strings.Contains(id, "btc") {
			ts.cache.InitTransactList(id)
		}
	}
	log.Printf("Number of connections --> %d", len(ts.cache.Pairs))
	return nil
}

func (ts *TickersService) ConcurrentTickDataCollection() {

	interrupt := make(chan os.Signal, 1)
	queues := make([]chan *models.Transaction, 40)
	for i := 0; i < 40; i++ {
		queue := make(chan *models.Transaction, 1)
		queues[i] = queue
	}
	signal.Notify(interrupt, os.Interrupt)
	var wg sync.WaitGroup
	j := 0

	for pair_key, _ := range ts.cache.Pairs {

		if j >= 40 {
			j = 0

		}
		pth := fmt.Sprintf("ws/" + pair_key + "@trade")
		u := url.URL{Scheme: "wss", Host: rootWsURI, Path: pth}
		wg.Add(1)
		go ts.exch.ConnectSocketAndReadTickData(u.String(), interrupt, queues[j], &wg)
		j++
	}

	log.Printf("Initialized %d websocket connections", j)

	for index, queue := range queues {
		go func(queue chan *models.Transaction, index int) {

			for {
				for transaction := range queue {

					if ts.writeToDB {
						log.Println("Writing to db -->", transaction)
						ts.store.Arr[index].InsertTransaction(transaction)

					} else {
						log.Println("Writing to cache -->", transaction)
						ts.cache.Insert(transaction)
					}
					ts.statsd.Client.Increment(".transacts.sqlinserts")
				}

			}

		}(queue, index)
	}

	go ts.statsd.ReportGoRoutines()
	wg.Wait()
}
