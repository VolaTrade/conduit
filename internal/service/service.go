package service

import (
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
	}

	CandlesService struct {
		cache *cache.CandlesCache
		db    *dynamo.DynamoSession
		rl    *limiter.RateLimiter
	}
)

func New(storage *dynamo.DynamoSession, ch *cache.CandlesCache) (*CandlesService, error) {

	var tempLimiter *limiter.RateLimiter
	var err error

	if tempLimiter, err = limiter.New(&limiter.Config{MaximumRequestPerInterval: 120, MinuteResetInterval: 1}); err != nil {
		return nil, err
	}
	return &CandlesService{cache: ch, db: storage, rl: tempLimiter}, nil
}
