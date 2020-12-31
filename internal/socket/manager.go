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
		logger                    *logger.Logger
		entry                     *models.CacheEntry
		txChannel                 chan *models.Transaction
		obChannel                 chan *models.OrderBookRow
		txSocketChan              chan bool
		transactionSocket         *ConduitSocket
		obFailSafeChan            chan bool
		transactionFailSafeSocket *ConduitSocket
		obSocketChan              chan bool
		orderBookFailSafeSocket   *ConduitSocket
		txFailSafeChan            chan bool
		orderBookSocket           *ConduitSocket
		kstats                    *stats.Stats
	}
)

//TODO add async startup for me
// TODO add health check functionality to me
//TODO add unit tests to me
// TOOD add wait group && context to me
func NewSocketManager(entry *models.CacheEntry, txChannel chan *models.Transaction, obChannel chan *models.OrderBookRow, statz *stats.Stats, logger *logger.Logger) *ConduitSocketManager {

	manager := &ConduitSocketManager{
		logger:                    logger,
		entry:                     entry,
		transactionSocket:         nil,
		txSocketChan:              make(chan bool),
		orderBookSocket:           nil,
		obSocketChan:              make(chan bool),
		transactionFailSafeSocket: nil,
		txFailSafeChan:            make(chan bool),
		orderBookFailSafeSocket:   nil,
		obFailSafeChan:            make(chan bool),
		txChannel:                 txChannel,
		obChannel:                 obChannel,
		kstats:                    statz,
	}

	return manager
}

func (csm *ConduitSocketManager) Run(ctx context.Context) {

	if err := csm.establishConnections(ctx); err != nil {
		panic(err)
	}
	go csm.consumeTransferTransactionMessage(ctx)
	go csm.consumeTransferOrderBookMessage(ctx)

}
func (csm *ConduitSocketManager) establishConnections(ctx context.Context) error {

	txSocket, err := NewSocket(ctx, csm.entry.TxUrl, csm.logger, csm.txSocketChan)

	if err != nil {
		csm.logger.Errorw("Socket failed to connect", "type", "transaction primary socket", "pair", csm.entry.Pair)
		return err
	}

	txFailSafeSocket, err := NewSocket(ctx, csm.entry.TxUrl, csm.logger, csm.txFailSafeChan)

	if err != nil {
		csm.logger.Errorw("Socket failed to connect", "type", "transaction failsafe socket", "pair", csm.entry.Pair)
		return err

	}

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
	csm.transactionSocket = txSocket
	csm.transactionFailSafeSocket = txFailSafeSocket

	return nil
}

func (csm *ConduitSocketManager) consumeTransferTransactionMessage(ctx context.Context) {
	csm.logger.Infow("Consuming and transferring messsage")
	mt := minuteTicker()
	for {

		message, err := csm.transactionSocket.readMessage()

		if err != nil {
			csm.logger.Errorw(err.Error(), "pair", csm.entry.Pair, "type", "transaction")
			csm.kstats.Increment("conduit.errors.socket_read.tx", 1.0)
			//---- tell socket to try reconnecting ---
			csm.txSocketChan <- false

			csm.logger.Infow("Performing context swap", "pair", csm.entry.Pair, "type", "transaction")

			// ---- swap ---
			csm.transactionSocket, csm.transactionFailSafeSocket = csm.transactionFailSafeSocket, csm.transactionSocket
			csm.txSocketChan, csm.txFailSafeChan = csm.txFailSafeChan, csm.txSocketChan
			// -------------
			continue

		}

		csm.kstats.Increment("conduit.socket_reads.tx", 1.0)

		var transaction *models.Transaction

		if transaction, err = models.UnmarshalTransactionJSON(message); err != nil {
			csm.logger.Errorw(err.Error(), "pair", csm.entry.Pair, "type", "transaction")
			csm.kstats.Increment("conduit.errors.json_unmarshal", 1.0)

		} else {
			csm.txChannel <- transaction
		}

		select {

		case <-mt.C:
			continue

		case <-ctx.Done():
			//csm.logger.Infow("received finish signal")
			return
		}

	}
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
