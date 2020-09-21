package service

import (
	"encoding/json"
	"net/http"

	"github.com/google/wire"
	"github.com/volatrade/candles/internal/binance"
	"github.com/volatrade/candles/internal/cache"
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
		db    *dynamo.CandlesDynamo
		exch  *binance.BinanceClient
	}
)

func New(storage *dynamo.CandlesDynamo, ch *cache.CandlesCache, cl *binance.BinanceClient) *CandlesService {

	return &CandlesService{cache: ch, db: storage, exch: cl}
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
		temp := val.(map[string]interface{}) //type casting
		id := temp["base"].(string) + temp["target"].(string)
		cs.cache.Pairs[id] = cache.InitializePair()
	}

	for key, _ := range cs.cache.Pairs {

		println("Initialized as key", key)
	}
	return nil
}

//Might make more sense to put this functionality into some API layer
func (cs *CandlesService) fetchCandleStick(pair string) error {

	//Append to five minute candle list here
	//ie cs.cache.InsertCandle(pairKey string, candle Candle) error {}
	return nil
}
