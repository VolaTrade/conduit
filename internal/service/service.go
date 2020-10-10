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
	"github.com/volatrade/candles/internal/dynamo"
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
		db    *dynamo.DynamoSession
		exch  *client.ApiClient
	}
)

func New(storage *dynamo.DynamoSession, ch *cache.CandlesCache, cl *client.ApiClient) *CandlesService {

	return &CandlesService{cache: ch, db: storage, exch: cl}
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

	for key, _ := range cs.cache.Pairs {

		println("Initialized as key", key)
	}
	return nil
}

func (cs *CandlesService) ConcurrentTickDataCollection() {

	interrupt := make(chan os.Signal, 1)
	queue := make(chan map[string]interface{}, 1)
	signal.Notify(interrupt, os.Interrupt)
	var wg sync.WaitGroup
	for pair_key, _ := range cs.cache.Pairs {
		pth := fmt.Sprintf("ws/" + pair_key + "@trade")
		u := url.URL{Scheme: "wss", Host: rootWsURI, Path: pth}
		wg.Add(1)
		go cs.exch.ConnectSocketAndReadTickData(u.String(), interrupt, queue, &wg)
	}

	go func() {

		for {
			for val := range queue {
				//call insert here
				log.Println("Val in queue: ", val)
			}

		}

	}()

	wg.Wait()
}
