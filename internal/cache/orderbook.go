package cache

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/volatrade/conduit/internal/models"
)

//getOrderBookUrlString builds orderbook websocket url from pair
func getOrderBookUrlString(pair string) string {
	innerPath := fmt.Sprintf("ws/" + strings.ToLower(pair) + "@depth10@1000ms")
	socketUrl := url.URL{Scheme: "wss", Host: BASE_SOCKET_URL, Path: innerPath}
	return socketUrl.String()
}

//GetAllOrderBookRows returns cache slice of OrderBookRow model struct
func (cc *ConduitCache) GetAllOrderBookRows() []*models.OrderBookRow {
	cc.obMux.RLock()
	defer cc.obMux.RUnlock()
	return cc.orderBookData
}

//InsertEntry takes pair, builds URLs, appends data to Entry model struct, then adds struct to cache
func (cc *ConduitCache) InsertOrderBookEntry(pair string) {

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
	cc.orderBookData = append(cc.orderBookData, obRow)
	cc.obMux.Unlock()
}

func (cc *ConduitCache) PurgeOrderBookRows() {
	cc.orderBookData = nil
}

//OrderBookRowsLength used for testing && debuging
func (cc *ConduitCache) OrderBookRowsLength() int {

	cc.obMux.RLock()
	defer cc.obMux.RUnlock()
	if cc.orderBookData != nil {
		return len(cc.orderBookData)
	}
	return 0
}
