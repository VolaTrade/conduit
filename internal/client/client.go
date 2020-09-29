package client

import (
	"github.com/google/wire"
	"github.com/volatrade/utilities/limiter"
)

var Module = wire.NewSet(
	New,
)

type Client interface {
	FetchFiveMinuteCandle(pair string) error
	GetActiveBinanceExchangePairs() ([]interface{}, error)
}

type ApiClient struct {
	rl *limiter.RateLimiter
}

func New() (*ApiClient, error) {
	var tempLimiter *limiter.RateLimiter
	var err error
	if tempLimiter, err = limiter.New(&limiter.Config{MaximumRequestPerInterval: 120, MinuteResetInterval: 1}); err != nil {
		return nil, err
	}

	return &ApiClient{rl: tempLimiter}, nil
}
