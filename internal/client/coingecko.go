package client

import (
	"encoding/json"
	"net/http"
)

type Symbol struct {
}

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
