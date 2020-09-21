package cache

import (
	"strconv"

	"github.com/google/wire"
)

var Module = wire.NewSet(
	New,
)

type (
	Cache interface {
	}

	Candle struct {
		Open      float64 `json:"open"`
		Close     float64 `json:"close"`
		High      float64 `json:"high"`
		Low       float64 `json:"low"`
		Timestamp string  `json:"timestamp"`
	}

	Pair struct {
		Five    []*Candle // 3
		Fifteen []*Candle // 2
		Thirty  []*Candle // 2
		Hour    []*Candle // 1
	}

	CandlesCache struct {
		Pairs map[string]*Pair
	}
)

/**
 * NewCandle does stuff
 */
func NewCandle(open string, close string, high string, low string, timestamp string) (*Candle, error) {
	output := &Candle{}

	value, err := strconv.ParseFloat(open, 64)
	if err != nil {
		return nil, err
	}
	output.Open = value

	value, err = strconv.ParseFloat(close, 64)
	if err != nil {
		return nil, err
	}
	output.Close = value

	value, err = strconv.ParseFloat(high, 64)
	if err != nil {
		return nil, err
	}
	output.High = value

	value, err = strconv.ParseFloat(low, 64)
	if err != nil {
		return nil, err
	}
	output.Low = value

	output.Timestamp = timestamp
	return output, nil

}

func InitializePair() *Pair {

	return &Pair{
		five:    make([]*Candle, 3),
		fifteen: make([]*Candle, 2),
		thirty:  make([]*Candle, 2),
		hour:    make([]*Candle, 2)}
}

func New() *CandlesCache {

	return &CandlesCache{Pairs: make(map[string]*Pair)}

}
