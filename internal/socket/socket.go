package socket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/volatrade/candles/internal/models"
)

type (
	Socket interface {
	}

	BinanceSocket struct {
		url         string
		dataChannel chan *models.Transaction
		connection  *websocket.Conn
	}
)

func NewSocket(urlString string, channel chan *models.Transaction) (*BinanceSocket, error) {
	println("establishing connection w/ socket for ->", urlString)
	conn, _, err := websocket.DefaultDialer.Dial(urlString, nil)

	if err != nil {
		return nil, err
	}
	println("conn established")
	socket := &BinanceSocket{url: urlString, connection: conn, dataChannel: channel}
	return socket, nil
}

func (bs *BinanceSocket) ReadMessage() ([]byte, error) {
	_, message, err := bs.connection.ReadMessage()

	if err != nil {
		return nil, err
	}
	return message, err
}

func (bs *BinanceSocket) readAndTransformTransaction() (*models.Transaction, error) {
	var json_message map[string]interface{}
	message, err := bs.ReadMessage()

	if err != nil {
		panic(err)
	}
	println("message -->", message)
	err = json.Unmarshal(message, &json_message)
	transaction, err := models.NewTransaction(json_message)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func ConsumeTransferMessage(socket *BinanceSocket) {
	errMax := 0
	for {

		transaction, err := socket.readAndTransformTransaction()

		if err != nil {
			//TODO ship error
			errMax++
			handleError(err, errMax)

		} else {

			socket.dataChannel <- transaction
		}
	}
}

func handleError(err error, errMax int) {
	if errMax == 3 {

		panic(err)
	}
	//Log here
}
