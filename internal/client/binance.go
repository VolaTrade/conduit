package client

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (ac *ApiClient) FetchFiveMinuteCandle(pair string) error {

	if !ac.rl.RequestsCanBeMade() {
		return errors.New("Maximum number of requests exceeded for interval")
	}

	endpoint := "https://api.binance.com/api/v1/klines?symbol=" + pair + "&interval=5m&limit=1"

	resp, err := http.Get(endpoint)

	ac.rl.IncrementRequestCount()

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(resp.Body)

	var result []interface{}
	if err := decoder.Decode(&result); err != nil {
		return err
	}

	//Implement the rest here
	return nil
}
