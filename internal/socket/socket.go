package socket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/volatrade/tickers/internal/models"
)

type (
	BinanceSocket struct {
		url         string
		dataChannel chan *models.Transaction
		connection  *websocket.Conn
	}
)

func NewSocket(urlString string, channel chan *models.Transaction) (*BinanceSocket, error) {

	socket := &BinanceSocket{url: urlString, connection: nil, dataChannel: channel}
	return socket, nil
}

func (bs *BinanceSocket) ReadMessage() ([]byte, error) {
	_, message, err := bs.connection.ReadMessage()

	if err != nil {
		println("MESSAGE -->", message)
		return nil, err
	}
	return message, err
}

func (bs *BinanceSocket) Connect() error {
	println("establishing connection w/ socket for ->", bs.url)
	conn, _, err := websocket.DefaultDialer.Dial(bs.url, nil)

	if err != nil {
		return err
	}
	println("conn established")
	bs.connection = conn
	return nil
}

func (bs *BinanceSocket) CloseConnection() error {
	return bs.connection.Close()
}

func readAndTransformTransaction(socket *BinanceSocket) (*models.Transaction, error) {
	var json_message map[string]interface{}
	message, err := socket.ReadMessage()

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(message, &json_message)
	transaction, err := models.NewTransaction(json_message)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func ConsumeTransferMessage(socket *BinanceSocket) {
	errMax := 0

	if err := socket.Connect(); err != nil {
		//TODO add handling policy
		panic(err)
	}
	for {

		transaction, err := readAndTransformTransaction(socket)
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
		println("Error MAX exceeded :p")
		panic(err)
	}
	//Log here
}
