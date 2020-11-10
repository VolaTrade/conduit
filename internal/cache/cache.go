package cache

import (
	"fmt"
	"net/url"

	"github.com/google/wire"
	"github.com/volatrade/candles/internal/models"
)

var Module = wire.NewSet(
	New,
)

const rootWsURI string = "stream.binance.com:9443"

type (
	Cache interface {
		InsertTransaction(transact *models.Transaction)
		PurgeTransactions()
		InsertPairUrl(pair string)
		GetPairUrl(index int) string
		GetTransaction(index int) *models.Transaction
		PairUrlsLength() int
		TransactionsLength() int
		GetAllTransactions() []*models.Transaction
	}
	TickersCache struct {
		pairUrls       []string
		transactions   []*models.Transaction
		transactLength int
	}
)

func getSocketUrlString(pair string) string {
	innerPath := fmt.Sprintf("ws/" + pair + "@trade")
	socketUrl := url.URL{Scheme: "wss", Host: rootWsURI, Path: innerPath}
	return socketUrl.String()
}
func (tc *TickersCache) PairUrlsLength() int {
	return len(tc.pairUrls)
}

func (tc *TickersCache) TransactionsLength() int {
	return tc.transactLength
}
func New() *TickersCache {
	return &TickersCache{
		pairUrls:       make([]string, 0),
		transactions:   make([]*models.Transaction, 0),
		transactLength: 0,
	}

}
func (tc *TickersCache) InsertTransaction(transact *models.Transaction) {
	tc.transactions = append(tc.transactions, transact)
	println(tc.transactions)
}

func (tc *TickersCache) InsertPairUrl(pair string) {
	tempUrl := getSocketUrlString(pair)
	tc.pairUrls = append(tc.pairUrls, tempUrl)
	println(tc.pairUrls)

}

func (tc *TickersCache) PurgeTransactions() {
	tc.transactions = nil
	tc.transactions = make([]*models.Transaction, 1)
	tc.transactLength = 0
}

func (tc *TickersCache) GetPairUrl(index int) string {
	return tc.pairUrls[index]
}

func (tc *TickersCache) GetTransaction(index int) *models.Transaction {
	return tc.transactions[index]
}
func (tc *TickersCache) GetAllTransactions() []*models.Transaction {
	return tc.transactions
}
