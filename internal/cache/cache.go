package cache

import (
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

const BASE_SOCKET_URL string = "stream.binance.com:9443"

type (
	Cache interface {
		InsertTransaction(transact *models.Transaction)
		InsertOrderBookRow(obRow *models.OrderBookRow)
		InsertTransactionUrl(pair string)
		InsertOrderBookUrl(pair string)
		GetTransactionUrl(index int) string
		GetOrderBookUrl(index int) string
		TransactionsLength() int
		OrderBookRowsLength() int
		GetAllTransactions() []*models.Transaction
		GetAllOrderBookRows() []*models.OrderBookRow
		UrlsLength() int
		Purge()
	}

	TickersCache struct {
		txUrls         []string
		obUrls         []string
		transactions   map[string][]*models.Transaction
		orderBookData  []*models.OrderBookRow
		transactLength int
		txMux          sync.Mutex
		obMux          sync.Mutex
	}
)

func New() *TickersCache {
	return &TickersCache{
		txUrls:         make([]string, 0),
		obUrls:         make([]string, 0),
		transactions:   make(map[string][]*models.Transaction, 0),
		orderBookData:  make([]*models.OrderBookRow, 0),
		transactLength: 0,
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

func (tc *TickersCache) InsertTransactionUrl(pair string) {
	tempUrl := getTransactionUrlString(pair)
	tc.txUrls = append(tc.txUrls, tempUrl)
	println(tc.txUrls)
}

func (tc *TickersCache) InsertOrderBookUrl(pair string) {
	tempUrl := getOrderBookUrlString(pair)
	tc.obUrls = append(tc.obUrls, tempUrl)

}

func (tc *TickersCache) UrlsLength() int {
	return len(tc.txUrls)
}

func (tc *TickersCache) TransactionsLength() int {
	return tc.transactLength
}

func (tc *TickersCache) OrderBookRowsLength() int {
	return len(tc.orderBookData)
}

func (tc *TickersCache) InsertTransaction(transact *models.Transaction) {
	tc.txMux.Lock()
	defer tc.txMux.Unlock()
	tc.transactions[transact.Pair] = append(tc.transactions[transact.Pair], transact)
	tc.transactLength++
}

func (tc *TickersCache) InsertOrderBookRow(obRow *models.OrderBookRow) {
	tc.obMux.Lock()
	defer tc.obMux.Unlock()
	tc.orderBookData = append(tc.orderBookData, obRow)
}

func (tc *TickersCache) Purge() {
	tc.transactions = nil
	tc.orderBookData = nil
	tc.transactLength = 0
}

func (tc *TickersCache) GetTransactionUrl(index int) string {
	return tc.txUrls[index]
}

func (tc *TickersCache) GetOrderBookUrl(index int) string {
	return tc.obUrls[index]
}

func (tc *TickersCache) GetAllTransactions() []*models.Transaction {
	tempTransacts := make([]*models.Transaction, tc.transactLength)

	for _, list := range tc.transactions {
		tempTransacts = append(tempTransacts, list...)
	}

	return tempTransacts
}

func (tc *TickersCache) GetAllOrderBookRows() []*models.OrderBookRow {
	return tc.orderBookData
}
