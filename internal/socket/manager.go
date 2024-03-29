package socket

import (
	"time"

	"context"

	"github.com/volatrade/conduit/internal/models"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/go-grafana-graphite-client"
)

type (
	ConduitSocketManager struct {
		logger                  *logger.Logger
		entry                   *models.CacheEntry
		obChannel               chan *models.OrderBookRow
		obFailSafeChan          chan bool
		obSocketChan            chan bool
		orderBookFailSafeSocket *ConduitSocket
		orderBookSocket         *ConduitSocket
		kstats                  stats.Stats
	}
)

//TODO add async startup for me
// TODO add health check functionality to me
//TODO add unit tests to me
// TOOD add wait group && context to me
func NewSocketManager(entry *models.CacheEntry, obChannel chan *models.OrderBookRow, statz stats.Stats, logger *logger.Logger) *ConduitSocketManager {

	manager := &ConduitSocketManager{
		logger:                  logger,
		entry:                   entry,
		orderBookSocket:         nil,
		obSocketChan:            make(chan bool),
		orderBookFailSafeSocket: nil,
		obFailSafeChan:          make(chan bool),
		obChannel:               obChannel,
		kstats:                  statz,
	}

	return manager
}

func (csm *ConduitSocketManager) Run(ctx context.Context) {

	if err := csm.establishConnections(ctx); err != nil {
		panic(err)
	}
	go csm.consumeTransferOrderBookMessage(ctx)

}
func (csm *ConduitSocketManager) establishConnections(ctx context.Context) error {

	obSocket, err := NewSocket(ctx, csm.entry.ObUrl, csm.logger, csm.obSocketChan)

	if err != nil {
		csm.logger.Errorw("Socket failed to connect", "type", "orderbook primary socket", "pair", csm.entry.Pair)
		return err
	}

	obFailSafeSocket, err := NewSocket(ctx, csm.entry.ObUrl, csm.logger, csm.obFailSafeChan)

	if err != nil {
		csm.logger.Errorw("Socket failed to connect", "type", "orderbook failsafe socket", "pair", csm.entry.Pair)
		return err
	}

	csm.orderBookSocket = obSocket
	csm.orderBookFailSafeSocket = obFailSafeSocket

	return nil
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

func (csm *ConduitSocketManager) consumeTransferOrderBookMessage(ctx context.Context) {
	csm.logger.Infow("Consuming and transferring orderbook message", "pair", csm.entry.Pair)
	mt := minuteTicker()
	defer mt.Stop()
	for {

		csm.logger.Infow("Reading order book message", "pair", csm.entry.Pair)
		message, err := csm.orderBookSocket.readMessage()

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
