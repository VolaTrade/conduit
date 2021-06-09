package candlerequest

import (
	"time"

	"context"

	"github.com/volatrade/conduit/internal/models"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/go-grafana-graphite-client"
)

type (
	CandleRequestManager struct {
		logger    *logger.Logger
		entry     *models.CacheEntry
		cdChannel chan *models.Kline
		kstats    stats.Stats
	}
)

//NewRequestManager ...
func NewCandleRequestManager(entry *models.CacheEntry, cdChannel chan *models.Kline, statz stats.Stats, logger *logger.Logger) *CandleRequestManager {

	manager := &CandleRequestManager{
		logger:    logger,
		entry:     entry,
		cdChannel: cdChannel,
		kstats:    statz,
	}

	return manager
}

func (crm *CandleRequestManager) Run(ctx context.Context) {

	go crm.consumeTransferCandlestickMessage(ctx)

}

func minuteTicker() *time.Ticker {
	c := make(chan time.Time, 1)
	t := &time.Ticker{C: c}
	go func() {
		for {
			n := time.Now()
			if n.Second() == 0 {
				c <- n
				time.Sleep(time.Second * 1)
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()
	return t
}

func (csm *CandleRequestManager) consumeTransferCandlestickMessage(ctx context.Context) {
	csm.logger.Infow("Consuming and transferring candlestick message", "pair", csm.entry.Pair)
	mt := minuteTicker()
	defer mt.Stop()
	for {

		csm.logger.Infow("Reading order book message", "pair", csm.entry.Pair)
		message, err := getRecentCandle()

		if err != nil {
			csm.logger.Errorw(err.Error(), "pair", csm.entry.Pair, "type", "orderbook")
			csm.kstats.Increment("errors.socket_read.ob", 1.0)

			//---- tell socket to try reconnecting ---
			csm.obSocketChan <- false

			// ---- context swap ---
			csm.logger.Infow("Performing context swap", "pair", csm.entry.Pair, "type", "orderbook")
			csm.orderBookSocket, csm.orderBookFailSafeSocket = csm.orderBookFailSafeSocket, csm.orderBookSocket
			csm.obSocketChan, csm.obFailSafeChan = csm.obFailSafeChan, csm.obSocketChan
			// -------------
			continue
		}
		csm.kstats.Increment("socket_reads.ob", 1.0)

		var orderBookRow *models.OrderBookRow

		if orderBookRow, err = models.UnmarshalOrderBookJSON(message, csm.entry.Pair); err != nil {
			csm.logger.Errorw(err.Error(), "pair", csm.entry.Pair, "type", "orderbook")
			csm.kstats.Increment("errors.json_unmarshal", 1.0)

		} else {
			csm.obChannel <- orderBookRow
		}

		time.Sleep(time.Second * 2)

		select {

		case <-mt.C:
			continue

		case <-ctx.Done():
			csm.logger.Infow("received finish signal")
			return
		}
	}

}
