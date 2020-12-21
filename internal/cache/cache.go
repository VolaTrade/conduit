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
		GetAllTransactions() []*models.Transaction
		InsertOrderBookRow(obRow *models.OrderBookRow)
		InsertTransaction(transact *models.Transaction)
		InsertEntry(pair string)
		GetEntries() []*models.CacheEntry
		OrderBookRowsLength() int
		Purge()
		TransactionsLength() int
	}

	ConduitCache struct {
		logger        *log.Logger
		entries       []*models.CacheEntry
		transactions  []*models.Transaction
		orderBookData []*models.OrderBookRow
		txMux         sync.Mutex
		obMux         sync.Mutex
	}
)

func New(logger *log.Logger) *ConduitCache {

	return &ConduitCache{
		logger:        logger,
		entries:       make([]*models.CacheEntry, 0),
		transactions:  make([]*models.Transaction, 0),
		orderBookData: make([]*models.OrderBookRow, 0),
	}

}
func getTransactionUrlString(pair string) string {
	innerPath := fmt.Sprintf("ws/" + strings.ToLower(pair) + "@trade")
	socketUrl := url.URL{Scheme: "wss", Host: BASE_SOCKET_URL, Path: innerPath}
	return socketUrl.String()
}
func getOrderBookUrlString(pair string) string {
	innerPath := fmt.Sprintf("ws/" + strings.ToLower(pair) + "@depth10@1000ms")
	socketUrl := url.URL{Scheme: "wss", Host: BASE_SOCKET_URL, Path: innerPath}
	return socketUrl.String()
}

func (tc *ConduitCache) GetAllTransactions() []*models.Transaction {
	return tc.transactions
}

func (tc *ConduitCache) GetAllOrderBookRows() []*models.OrderBookRow {
	return tc.orderBookData
}

func (cc *ConduitCache) InsertEntry(pair string) {

	entry := &models.CacheEntry{Pair: pair, TxUrl: getTransactionUrlString(pair), ObUrl: getOrderBookUrlString(pair)}
	cc.entries = append(cc.entries, entry)
}

func (cc *ConduitCache) InsertTransaction(transact *models.Transaction) {

	if transact == nil {
		return
	}

	cc.txMux.Lock()
	defer cc.txMux.Unlock()
	cc.transactions = append(cc.transactions, transact)
	cc.logger.Debugw("cache insertion", "type", "transaction", "current length", cc.OrderBookRowsLength())
}

func (cc *ConduitCache) InsertOrderBookRow(obRow *models.OrderBookRow) {
	if obRow == nil {
		return
	}

	cc.obMux.Lock()
	defer cc.obMux.Unlock()
	cc.orderBookData = append(cc.orderBookData, obRow)

	cc.logger.Debugw("cache insertion", "type", "orderbook snapshot", "current length", cc.OrderBookRowsLength())
}

func (cc *ConduitCache) Purge() {
	cc.transactions = nil
	cc.orderBookData = nil
}

func (cc *ConduitCache) GetEntries() []*models.CacheEntry {
	return cc.entries
}

//TransactionsLength used for testing && debuging
func (tc *ConduitCache) TransactionsLength() int {
	if tc.transactions != nil {
		return len(tc.transactions)
	}
	return 0
}

//OrderBookRowsLength used for testing && debuging
func (tc *ConduitCache) OrderBookRowsLength() int {

	if tc.orderBookData != nil {
		return len(tc.orderBookData)
	}
	return 0
}
