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
		candleRequest *CandleRequest
		logger        *logger.Logger
		entry         *models.CacheEntry
		cdChannel     chan *models.Kline
		cdRequestChan chan bool
		kstats        stats.Stats
	}
)

//NewRequestManager ...
func NewCandleRequestManager(entry *models.CacheEntry, cdChannel chan *models.Kline, statz stats.Stats, logger *logger.Logger) *CandleRequestManager {

	manager := &CandleRequestManager{
		logger:        logger,
		entry:         entry,
		cdChannel:     cdChannel,
		cdRequestChan: make(chan bool),
		kstats:        statz,
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

func (crm *CandleRequestManager) consumeTransferCandlestickMessage(ctx context.Context) {
	crm.logger.Infow("Consuming and transferring candlestick message", "pair", crm.entry.Pair)
	mt := minuteTicker()
	defer mt.Stop()
	for {

		crm.logger.Infow("Reading order book message", "pair", crm.entry.Pair)
		candlestick, err := crm.candleRequest.getRecentCandle(crm.entry.Pair)

		if err != nil {
			crm.logger.Errorw(err.Error(), "pair", crm.entry.Pair, "type", "candlestick")
			crm.kstats.Increment("errors.request_read.cd", 1.0)
			continue
		} else {
			crm.cdChannel <- candlestick
		}

		time.Sleep(time.Second * 2)

		select {

		case <-mt.C:
			continue

		case <-ctx.Done():
			crm.logger.Infow("received finish signal")
			return
		}
	}

}
