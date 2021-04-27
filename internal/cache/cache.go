//go:generate mockgen -package=mocks -destination=../mocks/cache.go github.com/volatrade/conduit/internal/cache Cache

package cache

import (
	"fmt"
	"net/url"
	"strings"
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
		GetAllOrderBookRows() []*models.OrderBookRow
		InsertOrderBookRow(obRow *models.OrderBookRow)
		InsertEntry(pair string)
		GetEntries() []*models.CacheEntry
		OrderBookRowsLength() int
		PurgeOrderBookRows()
	}

	ConduitCache struct {
		logger        *log.Logger
		entries       []*models.CacheEntry
		orderBookData []*models.OrderBookRow
		obMux         sync.Mutex
	}
)

//New ... constructor
func New(logger *log.Logger) *ConduitCache {

	return &ConduitCache{
		logger:        logger,
		entries:       make([]*models.CacheEntry, 0),
		orderBookData: make([]*models.OrderBookRow, 0),
	}

}

//getOrderBookUrlString builds orderbook websocket url from pair
func getOrderBookUrlString(pair string) string {
	innerPath := fmt.Sprintf("ws/" + strings.ToLower(pair) + "@depth10@1000ms")
	socketUrl := url.URL{Scheme: "wss", Host: BASE_SOCKET_URL, Path: innerPath}
	return socketUrl.String()
}

//GetAllOrderBookRows returns cache slice of OrderBookRow model struct
func (cc *ConduitCache) GetAllOrderBookRows() []*models.OrderBookRow {
	return cc.orderBookData
}

//InsertEntry takes pair, builds URLs, appends data to Entry model struct, then adds struct to cache
func (cc *ConduitCache) InsertEntry(pair string) {

	entry := &models.CacheEntry{Pair: pair, ObUrl: getOrderBookUrlString(pair)}
	cc.entries = append(cc.entries, entry)
}

//InsertOrderBookRow inserts OrderBookRow model struct to cache
func (cc *ConduitCache) InsertOrderBookRow(obRow *models.OrderBookRow) {
	if obRow == nil {
		cc.logger.Infow("Nil value passed in")
		return
	}

	cc.logger.Infow("cache insertion", "pair", obRow.Pair,
		"type", "orderbook snapshot", "cache length", cc.OrderBookRowsLength())
	cc.obMux.Lock()
	defer cc.obMux.Unlock()
	cc.orderBookData = append(cc.orderBookData, obRow)

}

func (cc *ConduitCache) PurgeOrderBookRows() {
	cc.orderBookData = nil

}

//GetEntries returns slice of CacheEntry struct
func (cc *ConduitCache) GetEntries() []*models.CacheEntry {
	return cc.entries
}

//OrderBookRowsLength used for testing && debuging
func (tc *ConduitCache) OrderBookRowsLength() int {

	if tc.orderBookData != nil {
		return len(tc.orderBookData)
	}
	return 0
}
