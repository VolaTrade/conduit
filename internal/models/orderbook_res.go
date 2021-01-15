package models

import (
	"encoding/json"
	"time"
)

type OrderBookRes struct {
	Id        int        `json:"lastUpdateId" db:"id"`
	Bids      [][]string `json:"bids" db:"bids"`
	Asks      [][]string `json:"asks" db:"asks"`
	Timestamp time.Time  `json:"timestamp" db:"timestamp"`
	Pair      string     `json:"pair" db:"pair"`
}

func UnmarshalDBOrderBookRow(obRow *OrderBookRow) (*OrderBookRes, error) {

	var obRes OrderBookRes
	err := json.Unmarshal(obRow.Bids, &obRes.Bids)
	if err != nil {
		return nil, err
	}

	err1 := json.Unmarshal(obRow.Asks, &obRes.Asks)
	if err1 != nil {
		return nil, err
	}
	obRes.Id = obRow.Id
	obRes.Timestamp = obRow.Timestamp
	obRes.Pair = obRow.Pair

	return &obRes, nil
}
