package cache

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/volatrade/conduit/internal/models"
)

//getCandleStickUrlString builds candlestick websocket url from pair
func getCandleStickUrlString(pair string) string {
	innerPath := fmt.Sprintf("ws/" + strings.ToLower(pair) + "@kline_1m")
	socketUrl := url.URL{Scheme: "wss", Host: BASE_SOCKET_URL, Path: innerPath}
	return socketUrl.String()
}

//GetAllCandleStickRows returns cache slice of CandleStickRow model struct
func (cc *ConduitCache) GetAllCandleStickRows() []*models.CandleStickRow {
	return cc.candleStickData
}

//InsertCandleStickRow inserts CandleStickRow model struct to cache
func (cc *ConduitCache) InsertCandleStickRow(cdRow *models.CandleStickRow) {
	if cdRow == nil {
		cc.logger.Infow("Nil value passed in for candle row")
		return
	}

	cc.logger.Infow("cache insertion", "type", "candlestick snapshot", "cache length", cc.CandleStickRowsLength())
	cc.cdMux.Lock()
	defer cc.cdMux.Unlock()
	cc.candleStickData = append(cc.candleStickData, cdRow)
}

//CandleStickRowsLength used for testing && debuging
func (tc *ConduitCache) CandleStickRowsLength() int {

	if tc.orderBookData != nil {
		return len(tc.candleStickData)
	}
	return 0
}
