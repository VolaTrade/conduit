package service

import (
	"github.com/google/wire"
	"github.com/volatrade/candles/internal/cache"
	"github.com/volatrade/candles/internal/dynamo"
)

var Module = wire.NewSet(
	New,
)

type (
	Service interface {
	}

	CandlesService struct {
		cache *cache.CandlesCache
		db    *dynamo.CandlesDynamo
	}
)

func New(storage *dynamo.CandlesDynamo, ch *cache.CandlesCache) *CandlesService {
	return &CandlesService{cache: ch, db: storage}
}
