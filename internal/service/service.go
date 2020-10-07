package service

import (
	"github.com/google/wire"
	"github.com/volatrade/candles/internal/cache"
	"github.com/volatrade/candles/internal/client"
	"github.com/volatrade/candles/internal/dynamo"
)

var Module = wire.NewSet(
	New,
)

type (
	Service interface {
		Init() error
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
		id := temp["base"].(string) + temp["target"].(string)
		cs.cache.Pairs[id] = cache.InitializePairData()
	}

	for key, _ := range cs.cache.Pairs {

		println("Initialized as key", key)
	}
	return nil
}
