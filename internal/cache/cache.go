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
		PurgeTransactions()
		InsertPairUrl(pair string)
		GetPairUrl(index int) string
		PairUrlsLength() int
		TransactionsLength() int
		GetAllTransactions() []*models.Transaction
	}
	TickersCache struct {
		pairUrls       []string
		transactions   map[string][]*models.Transaction
		transactLength int
		mux            sync.Mutex
	}
)

func getSocketUrlString(pair string) string {
	innerPath := fmt.Sprintf("ws/" + strings.ToLower(pair) + "@trade")
	socketUrl := url.URL{Scheme: "wss", Host: BASE_SOCKET_URL, Path: innerPath}
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
		transactions:   make(map[string][]*models.Transaction, 0),
		transactLength: 0,
	}

}
func (tc *TickersCache) InsertTransaction(transact *models.Transaction) {
	tc.mux.Lock()
	defer tc.mux.Unlock()
	tc.transactions[transact.Pair] = append(tc.transactions[transact.Pair], transact)
	tc.transactLength++
}

func (tc *TickersCache) InsertPairUrl(pair string) {
	tempUrl := getSocketUrlString(pair)
	tc.pairUrls = append(tc.pairUrls, tempUrl)
	println(tc.pairUrls)

}

func (tc *TickersCache) PurgeTransactions() {
	tc.transactions = nil
	tc.transactLength = 0
}

func (tc *TickersCache) GetPairUrl(index int) string {
	return tc.pairUrls[index]
}

func (tc *TickersCache) GetAllTransactions() []*models.Transaction {
	tempTransacts := make([]*models.Transaction, tc.transactLength)

	for _, list := range tc.transactions {
		tempTransacts = append(tempTransacts, list...)
	}

	return tempTransacts
}
