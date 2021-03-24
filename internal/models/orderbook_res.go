package models

import (
	"encoding/json"
	"time"
)

type OrderBookRes struct {
	Id              int           `json:"lastUpdateId" db:"id"`
	Bids            [][]string    `json:"bids" db:"bids"`
	Asks            [][]string    `json:"asks" db:"asks"`
	CreationTime    time.Time     `json:"time"`
	TransitDuration time.Duration `json:"duration"`
	Timestamp       string        `json:"timestamp" db:"timestamp"`
	Pair            string        `json:"pair" db:"pair"`
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
		Id:           obRow.Id,
		Bids:         obBids,
		Asks:         obAsks,
		CreationTime: obRow.CreationTime,
		Timestamp:    obRow.Timestamp,
		Pair:         obRow.Pair,
	}, nil
}

func (ob *OrderBookRes) UpdateTime(t time.Time) {
	ob.TransitDuration = time.Since(ob.CreationTime)
	ob.CreationTime = t
}
