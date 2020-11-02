package client

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/volatrade/candles/internal/models"
)

// GetActiveBinanceExchangePairs gets a list of all binance tradeable pairs
func (ac *ApiClient) GetActiveBinanceExchangePairs() ([]interface{}, error) {
	resp, err := http.Get("https://api.binance.com/api/v3/exchangeInfo")
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	var result map[string]interface{}

	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	dataPayLoad := result["symbols"].([]interface{})
	return dataPayLoad, nil
}

func (ac *ApiClient) ConnectSocketAndReadTickData(socketUrl string, interrupt chan os.Signal, queue chan *models.Transaction, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Connecting to %s", socketUrl)
	c, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Println("error")
		log.Printf("%e", err)
		wg.Done()
		ac.ConnectSocketAndReadTickData(socketUrl, interrupt, queue, wg)
		return
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("error")
				log.Printf("%e", err)
				wg.Done()
				c.Close()
				ac.ConnectSocketAndReadTickData(socketUrl, interrupt, queue, wg)
				return
			}
			var json_message map[string]interface{}

			err = json.Unmarshal(message, &json_message)
			if err != nil {
				//Log message
				log.Printf("%e", err)
			} else {
				transaction, err := models.NewTransaction(json_message)
				if err != nil {
					panic(err)
				}
				queue <- transaction
			}
		}
	}()
	for {

		select {
		case <-done:
			c.Close()
			return
		case <-interrupt:
			log.Println("Interrupt")
			//Cleanly close connection by sending a close message
			//then wait with timeout for the server to close the connection
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				panic(err)
			}
			return

		}
	}
}
