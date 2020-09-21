package client

import (
	"encoding/json"
	"net/http"
)

func (ac *ApiClient) GetActiveBinanceExchangePairs() ([]interface{}, error) {
	resp, err := http.Get("https://api.coingecko.com/api/v3/exchanges/binance/tickers")
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	var result map[string]interface{}

	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	dataPayLoad := result["tickers"].([]interface{})

	return dataPayLoad, nil
}
