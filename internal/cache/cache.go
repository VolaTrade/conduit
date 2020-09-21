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
		open      float64
		close     float64
		high      float64
		low       float64
		timestamp string
	}

	Pair struct {
		five    []*Candle // 3
		fifteen []*Candle // 2
		thirty  []*Candle // 2
		hour    []*Candle // 1
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
	output.open = value

	value, err = strconv.ParseFloat(close, 64)
	if err != nil {
		return nil, err
	}
	output.close = value

	value, err = strconv.ParseFloat(high, 64)
	if err != nil {
		return nil, err
	}
	output.high = value

	value, err = strconv.ParseFloat(low, 64)
	if err != nil {
		return nil, err
	}
	output.low = value

	output.timestamp = timestamp
	return output, nil

}

func InitializePair() *Pair {

	return &Pair{	
			five: make([]*Candle, 3), 
		     	fifteen: make([]*Candle, 2), 
			thirty: make([]*Candle, 2), 
			hour: make([]*Candle, 2)}
}

func New() *CandlesCache {

	return &CandlesCache{Pairs: make(map[string]*Pair)}

}

