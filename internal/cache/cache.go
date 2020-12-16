//go:generate mockgen -package=mocks -destination=../mocks/cache.go github.com/volatrade/conduit/internal/cache Cache

package cache

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/models"
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
		GetPair(index int) (string, error)
		GetTransactionOrderBookUrls(index int) (string, string, error)
		InsertOrderBookRow(obRow *models.OrderBookRow)
		InsertTransaction(transact *models.Transaction)
		InsertPair(pair string)
		OrderBookRowsLength() int
		PairsLength() int
		Purge()
		TransactionsLength() int
	}

	ConduitCache struct {
		pairs         []string
		transactions  []*models.Transaction
		orderBookData []*models.OrderBookRow
		txMux         sync.Mutex
		obMux         sync.Mutex
	}
)

func New() *ConduitCache {
	return &ConduitCache{
		pairs:         make([]string, 0),
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

func (tc *ConduitCache) GetTransactionOrderBookUrls(index int) (string, string, error) {

	if index < 0 || index >= len(tc.pairs) {
		return "", "", errors.New(OUT_OF_BOUNDS_ERROR)
	}
	return getTransactionUrlString(tc.pairs[index]), getOrderBookUrlString(tc.pairs[index]), nil
}

func (tc *ConduitCache) GetAllTransactions() []*models.Transaction {
	return tc.transactions
}

func (tc *ConduitCache) GetAllOrderBookRows() []*models.OrderBookRow {
	return tc.orderBookData
}

func (tc *ConduitCache) GetPair(index int) (string, error) {
	if index < 0 || index >= len(tc.pairs) {
		return "", errors.New(OUT_OF_BOUNDS_ERROR)
	}

	return tc.pairs[index], nil
}

func (tc *ConduitCache) TransactionsLength() int {
	if tc.transactions != nil {
		return len(tc.transactions)
	}
	return 0
}

func (tc *ConduitCache) OrderBookRowsLength() int {

	if tc.orderBookData != nil {
		return len(tc.orderBookData)
	}
	return 0
}

func (tc *ConduitCache) PairsLength() int {
	return len(tc.pairs)
}

func (tc *ConduitCache) InsertPair(pair string) {
	tc.pairs = append(tc.pairs, pair)

}

func (tc *ConduitCache) InsertTransaction(transact *models.Transaction) {

	if transact == nil {
		return
	}

	tc.txMux.Lock()
	defer tc.txMux.Unlock()
	tc.transactions = append(tc.transactions, transact)
}

func (tc *ConduitCache) InsertOrderBookRow(obRow *models.OrderBookRow) {
	log.Println("Inserting into cache", obRow)
	if obRow == nil {
		log.Println("Nil row case")
		return
	}

	tc.obMux.Lock()
	defer tc.obMux.Unlock()
	tc.orderBookData = append(tc.orderBookData, obRow)

	log.Println("Length", tc.OrderBookRowsLength())
}

func (tc *ConduitCache) Purge() {
	tc.transactions = nil
	tc.orderBookData = nil
	tc.pairs = nil
}
