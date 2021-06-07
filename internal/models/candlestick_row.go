package models

import (
	"encoding/json"
	"time"
)

type CandleStickRes struct {
	Id   int        `json:"lastUpdateId"`
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

type CandleStickRow struct {
	Id        int             `db:"id"`
	Bids      json.RawMessage `db:"bids"`
	Asks      json.RawMessage `db:"asks"`
	Timestamp time.Time       `db:"timestamp"`
	Pair      string          `db:"pair"`
}

func NewCandleStickRow(jsonResponse *CandleStickRes, pair string) (*CandleStickRow, error) {
	bids, err := json.Marshal(jsonResponse.Bids)
	if err != nil {
		return nil, err
	}

	asks, err := json.Marshal(jsonResponse.Asks)
	if err != nil {
		return nil, err
	}

	return &CandleStickRow{
		Id:        jsonResponse.Id,
		Timestamp: time.Now(),
		Pair:      pair,
		Asks:      asks,
		Bids:      bids,
	}, nil
}

// func UnmarshalCandleStickJSON(message []byte, pair string) (*OrderBookRow, error) {
// 	var jsonResponse CandleStickRes
// 	if err := json.Unmarshal(message, &jsonResponse); err != nil {
// 		return nil, err
// 	}
// 	ob, err := NewOrderBookRow(&jsonResponse, pair)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return ob, nil
// }
