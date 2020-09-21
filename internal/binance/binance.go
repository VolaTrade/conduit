package binance

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/wire"
	"github.com/volatrade/utilities/limiter"
)

var Module = wire.NewSet(
	New,
)

type Client interface{
	FetchFiveMinuteCandle(pair string) error

}

type BinanceClient struct {
	rl *limiter.RateLimiter
}

func New() (*BinanceClient, error) {
	var tempLimiter *limiter.RateLimiter
	var err error
	if tempLimiter, err = limiter.New(&limiter.Config{MaximumRequestPerInterval: 120, MinuteResetInterval: 1}); err != nil {
		return nil, err
	}

	return &BinanceClient{rl: tempLimiter}, nil
}

func (bc *BinanceClient) FetchFiveMinuteCandle(pair string) error {

	if !bc.rl.RequestsCanBeMade() {
		return errors.New("Maximum number of requests exceeded for interval")
	}

	endpoint := "https://api.binance.com/api/v1/klines?symbol=" + pair + "&interval=5m&limit=1"

	resp, err := http.Get(endpoint)

	bc.rl.IncrementRequestCount()

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(resp.Body)

	var result []interface{}
	if err := decoder.Decode(&result); err != nil {
		return err
	}

	//Implement the rest here
	return nil
}
