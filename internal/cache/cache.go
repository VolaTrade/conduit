package cache

import (
	"strconv"

	"github.com/google/wire"
	"github.com/volatrade/candles/internal/models"
)

var Module = wire.NewSet(
	New,
)

type Cache interface {
}

type CandlesCache struct {
	Pairs map[string]*models.PairData
}

/**
 * NewCandle does stuff
 */
func NewCandle(open string, close string, high string, low string) (*models.Candle, error) {
	output := &models.Candle{}

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

	return output, nil

}

func InitializePairData() *models.PairData {

	return &models.PairData{
		Five:    make([]*models.Candle, 3),
		Fifteen: make([]*models.Candle, 2),
		Thirty:  make([]*models.Candle, 2),
		Hour:    make([]*models.Candle, 2)}
}

func New() *CandlesCache {

	return &CandlesCache{Pairs: make(map[string]*models.PairData)}

}
