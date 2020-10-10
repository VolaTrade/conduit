package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

func (ac *ApiClient) FetchFiveMinuteCandle(pair string) error {

	if !ac.rl.RequestsCanBeMade() {
		return errors.New("Maximum number of requests exceeded for interval")
	}

	endpoint := "https://api.binance.com/api/v1/klines?symbol=" + pair + "&interval=5m&limit=1"

	resp, err := http.Get(endpoint)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Response error \n Status Code: %d \n Message: %s", resp.StatusCode, resp.Body))
	}

	ac.rl.IncrementRequestCount()
	decoder := json.NewDecoder(resp.Body)

	var result []interface{}
	if err := decoder.Decode(&result); err != nil {
		return err
	}

	//marshal data into candle struct
	return nil
}

func (ac *ApiClient) ConnectSocketAndReadTickData(u string, interrupt chan os.Signal, queue chan map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Connecting to %s", u)
	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		wg.Done()
		ac.ConnectSocketAndReadTickData(u, interrupt, queue, wg)
		return
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("%e", err)
				wg.Done()
				c.Close()
				ac.ConnectSocketAndReadTickData(u, interrupt, queue, wg)
				return
			}
			var json_message map[string]interface{}

			err = json.Unmarshal(message, &json_message)
			if err != nil {
				//Log message
				log.Printf("%e", err)
			} else {
				queue <- json_message
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
