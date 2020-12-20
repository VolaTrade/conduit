package requests

import (
	"encoding/json"
	"net/http"
	"strings"
)

// GetActiveBinanceExchangePairs gets a list of all binance tradeable pairs
func (ac *ConduitRequests) GetActiveBinanceExchangePairs() ([]string, error) {
	resp, err := http.Get("https://api.binance.com/api/v3/exchangeInfo")
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	var result map[string]interface{}

	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	rawData := result["symbols"].([]interface{})

	tradingPairs := make([]string, 0)

	for _, val := range rawData {
		temp := val.(map[string]interface{}) //type casting
		tradingPair := strings.ToLower(temp["symbol"].(string))
		tradingPairs = append(tradingPairs, tradingPair)
	}

	return tradingPairs, nil
}
