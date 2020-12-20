package socket

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/stats"
	logger "github.com/volatrade/currie-logs"
)

type (
	BinanceSocket struct {
		logger                *logger.Logger
		transactionUrl        string
		orderBookUrl          string
		Pair                  string
		TransactionChannel    chan *models.Transaction
		OrderBookChannel      chan *models.OrderBookRow
		transactionConnection *websocket.Conn
		orderBookConnection   *websocket.Conn
		statsd                *stats.StatsD
	}
)

func NewSocket(txUrl string, obUrl string, pair string, txChannel chan *models.Transaction,
	obChannel chan *models.OrderBookRow, statz *stats.StatsD, logger *logger.Logger) (*BinanceSocket, error) {

	socket := &BinanceSocket{
		transactionUrl:        txUrl,
		orderBookUrl:          obUrl,
		logger:                logger,
		Pair:                  pair,
		transactionConnection: nil,
		orderBookConnection:   nil,
		TransactionChannel:    txChannel,
		OrderBookChannel:      obChannel,
		statsd:                statz,
	}
	return socket, nil
}

func (bs *BinanceSocket) ReadMessage(messageType string) ([]byte, error) {

	var err error
	var message []byte
	if messageType == "TX" {
		_, message, err = bs.transactionConnection.ReadMessage()

	} else {
		_, message, err = bs.orderBookConnection.ReadMessage()
	}
	if err != nil {
		return nil, err
	}

	//println("Message received -->", message)

	//println("URl --->", bs.orderBookUrl)
	bs.statsd.Client.Increment(fmt.Sprintf("conduit.socket_reads.%s", bs.Pair))
	return message, err
}

func (bs *BinanceSocket) Connect() error {
	txConn, _, err := websocket.DefaultDialer.Dial(bs.transactionUrl, nil)

	if err != nil {
		return err
	}
	bs.logger.Infow("Connection established", "type", "transaction socket", "pair", bs.Pair)
	bs.transactionConnection = txConn

	obConn, _, err := websocket.DefaultDialer.Dial(bs.orderBookUrl, nil)

	if err != nil {
		return err
	}

	bs.logger.Infow("Connection established", "type", "orderbook socket", "pair", bs.Pair)
	bs.orderBookConnection = obConn

	return nil
}

func (bs *BinanceSocket) CloseConnections() (error, error) {
	return bs.transactionConnection.Close(), bs.orderBookConnection.Close()
}
