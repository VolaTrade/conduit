package models

import (
	"encoding/json"
	"log"
	"time"
)

type OrderBookRow struct {
	Id           int             `json:"lastUpdateId" db:"id"`
	Bids         json.RawMessage `json:"bids" db:"bids"`
	Asks         json.RawMessage `json:"asks" db:"asks"`
	CreationTime time.Time       `json:"time"`
	Timestamp    string          `json:"timestamp" db:"timestamp"`
	Pair         string          `json:"pair" db:"pair"`
}

func NewDBOrderBookRow(jsonResponse *OrderBookRes, pair string) (*OrderBookRow, error) {
	bids, err := json.Marshal(jsonResponse.Bids)
	if err != nil {
		return nil, err
	}

	asks, err := json.Marshal(jsonResponse.Asks)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	timestamp := t.Format("2006:01:02 15:04")

	if err != nil {
		log.Fatalf("Error forming custom timestamp: %s\n", err)
	}

	return &OrderBookRow{
		Id:           jsonResponse.Id,
		Bids:         bids,
		Asks:         asks,
		CreationTime: t,
		Timestamp:    timestamp,
		Pair:         pair,
	}, nil
}

func UnmarshalOrderBookJSON(message []byte, pair string) (*OrderBookRow, error) {
	var jsonResponse OrderBookRes
	if err := json.Unmarshal(message, &jsonResponse); err != nil {
		return nil, err
	}

	ob, err := NewDBOrderBookRow(&jsonResponse, pair)
	if err != nil {
		return nil, err
	}
	return ob, nil
}
