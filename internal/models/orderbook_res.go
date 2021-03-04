package models

import (
	"encoding/json"
	"time"
)

type OrderBookRes struct {
	Id        int        `json:"last_update_id" db:"id"`
	Bids      [][]string `json:"bids" db:"bids"`
	Asks      [][]string `json:"asks" db:"asks"`
	Time      time.Time  `json:"time"`
	Timestamp string     `json:"timestamp" db:"timestamp"`
	Pair      string     `json:"pair" db:"pair"`
}

func UnmarshalDBOrderBookRow(obRow *OrderBookRow) (*OrderBookRes, error) {

	var obBids [][]string
	err := json.Unmarshal(obRow.Bids, &obBids)
	if err != nil {
		return nil, err
	}

	var obAsks [][]string
	err1 := json.Unmarshal(obRow.Asks, &obAsks)
	if err1 != nil {
		return nil, err
	}

	return &OrderBookRes{
		Id:        obRow.Id,
		Bids:      obBids,
		Asks:      obAsks,
		Time:      obRow.Time,
		Timestamp: obRow.Timestamp,
		Pair:      obRow.Pair,
	}, nil
}
