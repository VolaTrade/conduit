package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/wire"
	"github.com/volatrade/candles/internal/cache"
	"github.com/volatrade/candles/internal/dynamo"
	"github.com/volatrade/utilities/limiter"
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
		db    *dynamo.CandlesDynamo
		rl    *limiter.RateLimiter
	}
)

func New(storage *dynamo.CandlesDynamo, ch *cache.CandlesCache) (*CandlesService, error) {

	var tempLimiter *limiter.RateLimiter
	var err error

	if tempLimiter, err = limiter.New(&limiter.Config{MaximumRequestPerInterval: 120, MinuteResetInterval: 1}); err != nil {
		return nil, err
	}
	return &CandlesService{cache: ch, db: storage, rl: tempLimiter}, nil
}

func (cs *CandlesService) Init() error {

	resp, err := http.Get("https://api.coingecko.com/api/v3/exchanges/binance/tickers")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(resp.Body)
	var result map[string]interface{}

	if err := decoder.Decode(&result); err != nil {
		return err
	}
	dataPayLoad := result["tickers"].([]interface{})
	for _, val := range dataPayLoad {
		temp := val.(map[string]interface{})
		coin_id := temp["base"].(string)
		coin_pair_id := temp["target"].(string)
		id := coin_id + coin_pair_id
		pair := cache.InitializePair()
		cs.cache.Pairs[id] = pair
	}

	for key, _ := range cs.cache.Pairs {

		println("Initialized as key", key)
	}
	return nil
}

//Might make more sense to put this functionality into some API layer
func (cs *CandlesService) fetchCandleStick(pair string) error {

	if !cs.rl.RequestsCanBeMade() {
		return errors.New("Maximum number of requests exceeded for interval")
	}

	endpoint := "https://api.binance.com/api/v1/klines?symbol=" + pair + "&interval=5m&limit=1"

	resp, err := http.Get(endpoint)

	cs.rl.IncrementRequestCount()

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(resp.Body)

	var result []interface{}
	if err := decoder.Decode(&result); err != nil {
		return err
	}

	//Append to five minute candle list here
	return nil
}
