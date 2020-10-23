package models

import (
	"log"
	"strconv"
	"time"
)

type Transaction struct {
	Pair      string    `db:"pair"`
	Price     float64   `db:"price"`
	IsMaker   bool      `db:"maker"`
	Timestamp time.Time `db:"timestamp"`
	Quantity  float64   `db:"quant"`
}

func NewTransaction(mapping map[string]interface{}) (*Transaction, error) {
	log.Println("Converting from ==>", mapping)
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

	quant, err := strconv.ParseFloat(quant_str, 64)
	if err != nil {
		return nil, err
	}
	return &Transaction{Timestamp: tm, Pair: pair, Price: price, Quantity: quant, IsMaker: maker}, nil
}
