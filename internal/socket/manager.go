package socket

import (
	"log"
	"time"

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

func NewSocketManager(entry *models.CacheEntry, txChannel chan *models.Transaction, obChannel chan *models.OrderBookRow, statz *stats.Stats, logger *logger.Logger) (*ConduitSocketManager, error) {

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

	if err := manager.EstablishConnections(); err != nil {
		return nil, err
	}

	go manager.consumeTransferTransactionMessage()
	go manager.consumeTransferOrderBookMessage()
	
	return manager, nil
}

func (bs *ConduitSocketManager) EstablishConnections() error {

	txSocket, err := NewSocket(bs.entry.TxUrl, bs.logger, bs.txSocketChan)

	if err != nil {
		bs.logger.Errorw("Socket failed to connect", "type", "transaction primary socket", "pair", bs.entry.Pair)
		return err
	}

	txFailSafeSocket, err := NewSocket(bs.entry.TxUrl, bs.logger, bs.txFailSafeChan)

	if err != nil {
		bs.logger.Errorw("Socket failed to connect", "type", "transaction failsafe socket", "pair", bs.entry.Pair)
		return err

	}

	obSocket, err := NewSocket(bs.entry.ObUrl, bs.logger, bs.obSocketChan)

	if err != nil {
		bs.logger.Errorw("Socket failed to connect", "type", "orderbook primary socket", "pair", bs.entry.Pair)
		return err
	}

	obFailSafeSocket, err := NewSocket(bs.entry.ObUrl, bs.logger, bs.obFailSafeChan)

	if err != nil {
		bs.logger.Errorw("Socket failed to connect", "type", "orderbook failsafe socket", "pair", bs.entry.Pair)
		return err
	}

	bs.orderBookSocket = obSocket
	bs.orderBookFailSafeSocket = obFailSafeSocket
	bs.transactionSocket = txSocket
	bs.transactionFailSafeSocket = txFailSafeSocket

	return nil
}

func (csm *ConduitSocketManager) consumeTransferTransactionMessage() {
	csm.logger.Infow("Consuming and transferring messsage")

	hits := 0
	time_start := time.Now()
	for {

		if time_elapsed := time.Since(time_start); time_elapsed <= time.Duration(time.Millisecond*200) {

			continue
		}

		message, err := csm.transactionSocket.readMessage()
		hits++

		if err != nil {
			csm.logger.Errorw(err.Error())
			csm.kstats.Increment("conduit.errors.socket_read.tx", 1.0)
			//---- tell socket to try reconnecting ---
			csm.txSocketChan <- false

			// ---- swap ---
			csm.transactionSocket, csm.transactionFailSafeSocket = csm.transactionFailSafeSocket, csm.transactionSocket
			csm.txSocketChan, csm.txFailSafeChan = csm.txFailSafeChan, csm.txSocketChan
			// -------------

			continue

		}

		csm.kstats.Increment("conduit.socket_reads.tx", 1.0)

		var transaction *models.Transaction

		if transaction, err = models.UnmarshalTransactionJSON(message); err != nil {
			csm.logger.Errorw(err.Error())
			csm.kstats.Increment("conduit.errors.json_unmarshal", 1.0)

		} else {
			csm.txChannel <- transaction
		}

	}
}

func (csm *ConduitSocketManager) consumeTransferOrderBookMessage() {
	log.Println("Consuming and transferring messsage")

	prev_min := time.Now().Minute() - 1
	for {

		curr_min := time.Now().Minute()

		if prev_min == curr_min {
			continue
		}
		message, err := csm.orderBookSocket.readMessage()
		prev_min = curr_min
		if err != nil {
			//handle me
			csm.logger.Errorw(err.Error(), csm.entry.Pair)
			csm.kstats.Increment("conduit.errors.socket_read.ob", 1.0)

			//---- tell socket to try reconnecting ---
			csm.obSocketChan <- false

			// ---- context swap ---
			csm.orderBookSocket, csm.orderBookFailSafeSocket = csm.orderBookFailSafeSocket, csm.orderBookSocket
			csm.obSocketChan, csm.obFailSafeChan = csm.obFailSafeChan, csm.obSocketChan
			// -------------
			continue
		}

		csm.kstats.Increment("conduit.socket_reads.ob", 1.0)

		var orderBookRow *models.OrderBookRow

		if orderBookRow, err = models.UnmarshalOrderBookJSON(message, csm.entry.Pair); err != nil {
			log.Println(err.Error(), csm.entry.Pair)
			csm.kstats.Increment("conduit.errors.json_unmarshal", 1.0)

		} else {
			csm.obChannel <- orderBookRow
		}
	}

}
