package socket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/volatrade/conduit/internal/models"
	stats "github.com/volatrade/k-stats"
)

type (
	BinanceSocket struct {
		transactionUrl        string
		orderBookUrl          string
		Pair                  string
		TransactionChannel    chan *models.Transaction
		OrderBookChannel      chan *models.OrderBookRow
		transactionConnection *websocket.Conn
		orderBookConnection   *websocket.Conn
		kstats                *stats.Stats
	}
)

func NewSocket(txUrl string, obUrl string, pair string, txChannel chan *models.Transaction, obChannel chan *models.OrderBookRow, statz *stats.Stats) (*BinanceSocket, error) {

	socket := &BinanceSocket{
		transactionUrl:        txUrl,
		orderBookUrl:          obUrl,
		Pair:                  pair,
		transactionConnection: nil,
		orderBookConnection:   nil,
		TransactionChannel:    txChannel,
		OrderBookChannel:      obChannel,
		kstats:                statz,
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
	bs.kstats.Increment(fmt.Sprintf("conduit.socket_reads.%s", bs.Pair), 1.0)
	return message, err
}

func (bs *BinanceSocket) Connect() error {
	log.Println("establishing connection w/ socket for ->", bs.transactionUrl)
	txConn, _, err := websocket.DefaultDialer.Dial(bs.transactionUrl, nil)

	if err != nil {
		return err
	}
	log.Println("conn established for transaction socket")
	bs.transactionConnection = txConn

	obConn, _, err := websocket.DefaultDialer.Dial(bs.orderBookUrl, nil)

	if err != nil {
		return err
	}

	log.Println("conn established for order book socket")
	bs.orderBookConnection = obConn

	return nil
}

func (bs *BinanceSocket) CloseConnections() (error, error) {
	return bs.transactionConnection.Close(), bs.orderBookConnection.Close()
}
