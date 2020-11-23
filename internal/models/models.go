package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Transaction struct {
	Id        int64     `db:"id"`
	Pair      string    `db:"pair"`
	Price     float64   `db:"price"`
	IsMaker   bool      `db:"maker"`
	Timestamp time.Time `db:"timestamp"`
	Quantity  float64   `db:"quant"`
}

type OrderBookRes struct {
	Id   int        `json:"lastUpdateId" db:"id"`
	Bids [][]string `json:"bids" db:"bids"`
	Asks [][]string `json:"asks" db:"asks"`
}

type OrderBookRow struct {
	Id        int         `db:"id"`
	Bids      [][]float64 `db:"bids"`
	Asks      [][]float64 `db:"asks"`
	Timestamp time.Time   `db:"timestamp"`
	Pair      string      `db:"pair"`
}

func NewOrderBookRow(jsonResponse *OrderBookRes, pair string) (*OrderBookRow, error) {
	bids, err := Str2FloatSlice(jsonResponse.Bids)
	if err != nil {
		return nil, err
	}
	asks, err := Str2FloatSlice(jsonResponse.Asks)
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

func NewTransaction(mapping map[string]interface{}) (*Transaction, error) {
	i := int64(mapping["T"].(float64)) / 1000
	tm := time.Unix(i, 0)
	price_str := mapping["p"].(string)
	price, err := strconv.ParseFloat(price_str, 64)
	if err != nil {
		return nil, err
	}
	pair := mapping["s"].(string)
	maker := mapping["m"].(bool)
	quant_str := mapping["q"].(string)
	id := int64(mapping["a"].(float64))

	quant, err := strconv.ParseFloat(quant_str, 64)
	if err != nil {
		return nil, err
	}
	return &Transaction{Id: id, Timestamp: tm, Pair: pair, Price: price, Quantity: quant, IsMaker: maker}, nil
}

func UnmarshalTransactionJSON(message []byte) (*Transaction, error) {
	var json_message map[string]interface{}

	if err := json.Unmarshal(message, &json_message); err != nil {
		return nil, err
	}

	return NewTransaction(json_message)

}

func Str2FloatSlice(sl [][]string) ([][]float64, error) {

	n := make([][]float64, len(sl))
	for i := range sl {
		n[i] = make([]float64, 0)

		for j := 0; j < len(sl[i]); j++ {
			f64, err := strconv.ParseFloat(sl[i][j], 64)
			if err != nil {
				return nil, err
			}
			n[i] = append(n[i], f64)
		}

	}
	fmt.Println(n)
	return n, nil
}
