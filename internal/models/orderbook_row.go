package models

import (
	"encoding/json"
	"time"
)

type OrderBookRow struct {
	Id        int        `json:"lastUpdateId"`
	Bids      [][]string `json:"bids"`
	Asks      [][]string `json:"asks"`
	Timestamp time.Time  `json:"timestamp"`
	Pair      string     `json:"pair"`
}

func UnmarshalOrderBookJSON(message []byte, pair string) (*OrderBookRow, error) {
	var jsonResponse OrderBookRow
	if err := json.Unmarshal(message, &jsonResponse); err != nil {
		return nil, err
	}

	jsonResponse.Timestamp = time.Now()
	jsonResponse.Pair = pair

	return &jsonResponse, nil
}
