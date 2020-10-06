package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (ac *ApiClient) FetchFiveMinuteCandle(pair string) error {

	if !ac.rl.RequestsCanBeMade() {
		return errors.New("Maximum number of requests exceeded for interval")
	}

	endpoint := "https://api.binance.com/api/v1/klines?symbol=" + pair + "&interval=5m&limit=1"

	resp, err := http.Get(endpoint)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

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
