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

//TODO add async startup for me
// TODO add health check functionality to me
//TODO add unit tests to me
// TOOD add wait group && context to me
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

func (csm *ConduitSocketManager) EstablishConnections() error {

	txSocket, err := NewSocket(csm.entry.TxUrl, csm.logger, csm.txSocketChan)

	if err != nil {
		csm.logger.Errorw("Socket failed to connect", "type", "transaction primary socket", "pair", csm.entry.Pair)
		return err
	}

	txFailSafeSocket, err := NewSocket(csm.entry.TxUrl, csm.logger, csm.txFailSafeChan)

	if err != nil {
		csm.logger.Errorw("Socket failed to connect", "type", "transaction failsafe socket", "pair", csm.entry.Pair)
		return err

	}

	obSocket, err := NewSocket(csm.entry.ObUrl, csm.logger, csm.obSocketChan)

	if err != nil {
		csm.logger.Errorw("Socket failed to connect", "type", "orderbook primary socket", "pair", csm.entry.Pair)
		return err
	}

	obFailSafeSocket, err := NewSocket(csm.entry.ObUrl, csm.logger, csm.obFailSafeChan)

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

func (csm *ConduitSocketManager) consumeTransferTransactionMessage() {
	csm.logger.Infow("Consuming and transferring messsage")

	ticker := time.NewTicker(time.Millisecond * 200)

	for {

		message, err := csm.transactionSocket.readMessage()

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

		select {

		case <-ticker.C:
			continue
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
			}
			time.Sleep(time.Second)
		}
	}()
	return t
}

func (csm *ConduitSocketManager) consumeTransferOrderBookMessage() {
	log.Println("Consuming and transferring messsage")

	for {

		message, err := csm.orderBookSocket.readMessage()

		if err != nil {
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

		time.Sleep(time.Second * 2)
		select {

		case <-minuteTicker().C:
			continue
		}
	}

}
