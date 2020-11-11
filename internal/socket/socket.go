package socket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/volatrade/tickers/internal/models"
	"github.com/volatrade/tickers/internal/stats"
)

type (
	BinanceSocket struct {
		url         string
		pair        string
		DataChannel chan *models.Transaction
		connection  *websocket.Conn
		statsd      *stats.StatsD
	}
)

func NewSocket(urlString string, pair string, channel chan *models.Transaction, statz *stats.StatsD) (*BinanceSocket, error) {

	socket := &BinanceSocket{url: urlString, pair: pair, connection: nil, DataChannel: channel, statsd: statz}
	return socket, nil
}

func (bs *BinanceSocket) ReadMessage() ([]byte, error) {
	_, message, err := bs.connection.ReadMessage()

	if err != nil {
		log.Println("message from error ->", message)
		return nil, err
	}
	bs.statsd.Client.Increment(fmt.Sprintf("tickers.socket_reads.%s", bs.pair))
	return message, err
}

func (bs *BinanceSocket) Connect() error {
	log.Println("establishing connection w/ socket for ->", bs.url)
	conn, _, err := websocket.DefaultDialer.Dial(bs.url, nil)

	if err != nil {
		return err
	}
	log.Println("conn established")
	bs.connection = conn
	return nil
}

func (bs *BinanceSocket) CloseConnection() error {
	return bs.connection.Close()
}
