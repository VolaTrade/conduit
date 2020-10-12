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
	"github.com/volatrade/candles/internal/storage"
)

const rootWsURI string = "stream.binance.com:9443"

var Module = wire.NewSet(
	New,
)

type (
	Service interface {
		Init() error
		ConcurrentTickDataCollection()
	}

	CandlesService struct {
		cache *cache.CandlesCache
		store *storage.ConnectionArray
		exch  *client.ApiClient
	}
)

func New(arr *storage.ConnectionArray, ch *cache.CandlesCache, cl *client.ApiClient) *CandlesService {

	return &CandlesService{cache: ch, store: arr, exch: cl}
}

func (cs *CandlesService) Init() error {

	tradingCryptosList, err := cs.exch.GetActiveBinanceExchangePairs()
	if err != nil {
		return err
	}

	for _, val := range tradingCryptosList {
		temp := val.(map[string]interface{}) //type casting
		id := strings.ToLower(temp["base"].(string) + temp["target"].(string))
		cs.cache.Pairs[id] = cache.InitializePairData()
	}
	log.Printf("Number of connections --> %d", len(cs.cache.Pairs))
	return nil
}

func (cs *CandlesService) ConcurrentTickDataCollection() {

	interrupt := make(chan os.Signal, 1)
	queues := make([]chan map[string]interface{}, 40)
	for i := 0; i < 40; i++ {
		queue := make(chan map[string]interface{}, 1)
		queues[i] = queue
	}
	signal.Notify(interrupt, os.Interrupt)
	var wg sync.WaitGroup
	j := 0

	for pair_key, _ := range cs.cache.Pairs {

		if j >= 40 {
			j = 0

		}
		pth := fmt.Sprintf("ws/" + pair_key + "@trade")
		u := url.URL{Scheme: "wss", Host: rootWsURI, Path: pth}
		wg.Add(1)
		go cs.exch.ConnectSocketAndReadTickData(u.String(), interrupt, queues[j], &wg)
		j++
	}

	log.Printf("Initialized %d connections", j)
	for index, queue := range queues {
		go func(queue chan map[string]interface{}, index int) {

			for {
				for val := range queue {

					cs.store.Arr[index].InsertTransaction(val)
					log.Printf("Val in queue:  %s @ queue #%d w/ queue length -> %d", val, index, len(queue))
				}

			}

		}(queue, index)
	}

	wg.Wait()
}
