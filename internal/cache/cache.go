//go:generate mockgen -package=mocks -destination=../mocks/cache.go github.com/volatrade/conduit/internal/cache Cache

package cache

import (
	"sync"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/models"
	log "github.com/volatrade/currie-logs"
)

var Module = wire.NewSet(
	New,
)

const (
	BASE_SOCKET_URL     string = "stream.binance.com:9443"
	OUT_OF_BOUNDS_ERROR string = "Index does not exist for pair slice"
)

type (
	Cache interface {
		GetAllCandleStickRows() []*models.CandleStickRow
		GetAllOrderBookRows() []*models.OrderBookRow
		InsertCandleStickRow(cdRow *models.CandleStickRow)
		InsertOrderBookRow(obRow *models.OrderBookRow)
		InsertOrderBookEntry(pair string)
		GetEntries() []*models.CacheEntry
		CandleStickRowsLength() int
		OrderBookRowsLength() int
		PurgeOrderBookRows()
	}

	ConduitCache struct {
		logger          *log.Logger
		entries         []*models.CacheEntry
		orderBookData   []*models.OrderBookRow
		candleStickData []*models.CandleStickRow
		obMux           *sync.RWMutex
		cdMux           *sync.RWMutex
	}
)

//New ... constructor
func New(logger *log.Logger) *ConduitCache {

	return &ConduitCache{
		logger:          logger,
		entries:         make([]*models.CacheEntry, 0),
		obMux:           &sync.RWMutex{},
		cdMux:           &sync.RWMutex{},
		orderBookData:   make([]*models.OrderBookRow, 0),
		candleStickData: make([]*models.CandleStickRow, 0),
	}

}

//GetEntries returns slice of CacheEntry struct
func (cc *ConduitCache) GetEntries() []*models.CacheEntry {
	return cc.entries
}
