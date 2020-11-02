package cache

import (
	"github.com/google/wire"
	"github.com/volatrade/candles/internal/models"
)

var Module = wire.NewSet(
	New,
)

type Cache interface {
	InitTransactList(pair string)
	Insert(transact *models.Transaction)
	Purge()
}

type TickersCache struct {
	Pairs map[string][]*models.Transaction
}

func New() *TickersCache {
	return &TickersCache{Pairs: make(map[string][]*models.Transaction)}

}
func (tc *TickersCache) Insert(transact *models.Transaction) {
	tc.Pairs[transact.Pair] = append(tc.Pairs[transact.Pair], transact)
}

func (tc *TickersCache) InitTransactList(pair string) {
	tc.Pairs[pair] = make([]*models.Transaction, 10)
}

func (tc *TickersCache) Purge() {
	for key, _ := range tc.Pairs {
		tc.Pairs[key] = nil
	}
}
