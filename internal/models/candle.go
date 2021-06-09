package models

import (
	"encoding/json"
	"strconv"
	"time"
)

type Kline struct {
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
	Pair      string    `db:"pair"`
}

func UnmarshallBytesToLatestKline(bits []byte, pair string) (*Kline, error) {
	var rawInterfaces [][]interface{}
	var err error
	var kline Kline

	kline.Pair = pair

	if err := json.Unmarshal(bits, &rawInterfaces); err != nil {
		return nil, err
	}
	kline.Timestamp = time.Unix(int64(rawInterfaces[0][0].(float64)), 0)

	kline.Open, err = convertToFloat64(rawInterfaces[0][1].(string))
	if err != nil {
		return nil, err
	}
	kline.High, err = convertToFloat64(rawInterfaces[0][2].(string))
	if err != nil {
		return nil, err
	}
	kline.Low, err = convertToFloat64(rawInterfaces[0][3].(string))
	if err != nil {
		return nil, err
	}
	kline.Close, err = convertToFloat64(rawInterfaces[0][4].(string))
	if err != nil {
		return nil, err
	}
	kline.Volume, err = convertToFloat64(rawInterfaces[0][5].(string))
	if err != nil {
		return nil, err
	}
	return &kline, nil
}

func convertToFloat64(str string) (float64, error) {
	float, err := strconv.ParseFloat(str, 64)

	if err != nil {
		return 0, err
	}

	return float, nil
}
