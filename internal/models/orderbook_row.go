package models

import (
	"encoding/json"
	"time"
)

type OrderBookRes struct {
	Id   int        `json:"lastUpdateId"`
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

type OrderBookRow struct {
	Id        int             `db:"id"`
	Bids      json.RawMessage `db:"bids"`
	Asks      json.RawMessage `db:"asks"`
	Timestamp time.Time       `db:"timestamp"`
	Pair      string          `db:"pair"`
}

func NewOrderBookRow(jsonResponse *OrderBookRes, pair string) (*OrderBookRow, error) {
	bids, err := json.Marshal(jsonResponse.Bids)
	if err != nil {
		return nil, err
	}

	asks, err := json.Marshal(jsonResponse.Asks)
	if err != nil {
		return nil, err
	}

	return &OrderBookRow{
		Id:        jsonResponse.Id,
		Timestamp: time.Now(),
		Pair:      pair,
		Asks:      asks,
		Bids:      bids,
	}, nil
}

func UnmarshalOrderBookJSON(message []byte, pair string) (*OrderBookRow, error) {
	var jsonResponse OrderBookRes
	if err := json.Unmarshal(message, &jsonResponse); err != nil {
		return nil, err
	}
	ob, err := NewOrderBookRow(&jsonResponse, pair)
	if err != nil {
		return nil, err
	}

	return ob, nil
}
