package models

import (
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
