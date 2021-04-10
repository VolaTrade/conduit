package socket

import (
	"time"

	"context"

	"github.com/volatrade/conduit/internal/models"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
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

// type (
// 	CortexSocketManager struct {
// 		logger         *logger.Logger
// 		cortexEntry    *models.CortexEntry
// 		pedersonSocket *ConduitSocket
// 		kstats         stats.Stats
// 	}
// )

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
	csm.logger.Infow("Consuming and transferring orderbook message")
	mt := minuteTicker()
	defer mt.Stop()
	for {

		csm.logger.Infow("Reading order book message", "pair", csm.entry.Pair)
		message, err := csm.orderBookSocket.readMessage()

		csm.logger.Infow("Message read.. checking for error", "pair", csm.entry.Pair)
		if err != nil {
			csm.logger.Errorw(err.Error(), "pair", csm.entry.Pair, "type", "orderbook")
			csm.kstats.Increment("conduit.errors.socket_read.ob", 1.0)

			//---- tell socket to try reconnecting ---
			csm.obSocketChan <- false

			// ---- context swap ---
			csm.logger.Infow("Performing context swap", "pair", csm.entry.Pair, "type", "orderbook")
			csm.orderBookSocket, csm.orderBookFailSafeSocket = csm.orderBookFailSafeSocket, csm.orderBookSocket
			csm.obSocketChan, csm.obFailSafeChan = csm.obFailSafeChan, csm.obSocketChan
			// -------------
			continue
		}

		timeNow := time.Now()
		csm.logger.Infow("Incrementing", "pair", csm.entry.Pair)
		csm.kstats.Increment("conduit.socket_reads.ob", 1.0)
		csm.logger.Infow("Increment complete", "time", time.Since(timeNow).String())

		var orderBookRow *models.OrderBookRow

		if orderBookRow, err = models.UnmarshalOrderBookJSON(message, csm.entry.Pair); err != nil {
			csm.logger.Errorw(err.Error(), "pair", csm.entry.Pair, "type", "orderbook")
			csm.kstats.Increment("conduit.errors.json_unmarshal", 1.0)

		} else {
			csm.logger.Infow("Sending orderbook data to channel", "pair", csm.entry.Pair, "type", "orderbook")
			csm.obChannel <- orderBookRow
			csm.logger.Infow("Data sent through channel", "pair", csm.entry.Pair, "type", "orderbook")
		}

		time.Sleep(time.Second * 2)

		csm.logger.Infow("Going into select", "pair", csm.entry.Pair, "type", "orderbook")
		select {

		case <-mt.C:
			csm.logger.Infow("Got ticker signal for orderbook data", "pair", csm.entry.Pair)
			continue

		case <-ctx.Done():
			csm.logger.Infow("received finish signal")
			return
		}
	}

}

// func NewPedersonSocketManager(entry *models.CortexEntry, statz stats.Stats, logger *logger.Logger) *CortexSocketManager {

// 	manager := &CortexSocketManager{
// 		logger:         logger,
// 		cortexEntry:    entry,
// 		pedersonSocket: nil,
// 		kstats:         statz,
// 	}

// 	return manager
// }
