package cache

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/google/wire"
	"github.com/volatrade/tickers/internal/models"
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

	TickersCache struct {
		pairs         []string
		transactions  []*models.Transaction
		orderBookData []*models.OrderBookRow
		txMux         sync.Mutex
		obMux         sync.Mutex
	}
)

func New() *TickersCache {
	return &TickersCache{
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
	innerPath := fmt.Sprintf("ws/" + strings.ToLower(pair) + "@depth10@100ms")
	socketUrl := url.URL{Scheme: "wss", Host: BASE_SOCKET_URL, Path: innerPath}
	return socketUrl.String()
}

func (tc *TickersCache) GetTransactionOrderBookUrls(index int) (string, string, error) {

	if index < 0 || index >= len(tc.pairs) {
		return "", "", errors.New(OUT_OF_BOUNDS_ERROR)
	}
	return getTransactionUrlString(tc.pairs[index]), getOrderBookUrlString(tc.pairs[index]), nil
}

func (tc *TickersCache) GetAllTransactions() []*models.Transaction {
	return tc.transactions
}

func (tc *TickersCache) GetAllOrderBookRows() []*models.OrderBookRow {
	return tc.orderBookData
}

func (tc *TickersCache) GetPair(index int) (string, error) {
	if index < 0 || index >= len(tc.pairs) {
		return "", errors.New(OUT_OF_BOUNDS_ERROR)
	}

	return tc.pairs[index], nil
}

func (tc *TickersCache) TransactionsLength() int {
	if tc.transactions != nil {
		return len(tc.transactions)
	}
	return 0
}

func (tc *TickersCache) OrderBookRowsLength() int {

	if tc.orderBookData != nil {
		return len(tc.orderBookData)
	}
	return 0
}

func (tc *TickersCache) PairsLength() int {
	return len(tc.pairs)
}

func (tc *TickersCache) InsertPair(pair string) {
	tc.pairs = append(tc.pairs, pair)

}

func (tc *TickersCache) InsertTransaction(transact *models.Transaction) {

	if transact == nil {
		return
	}

	tc.txMux.Lock()
	defer tc.txMux.Unlock()
	tc.transactions = append(tc.transactions, transact)
}

func (tc *TickersCache) InsertOrderBookRow(obRow *models.OrderBookRow) {

	if obRow == nil {
		return
	}

	tc.obMux.Lock()
	defer tc.obMux.Unlock()
	tc.orderBookData = append(tc.orderBookData, obRow)
}

func (tc *TickersCache) Purge() {
	tc.transactions = nil
	tc.orderBookData = nil
	tc.pairs = nil
}
