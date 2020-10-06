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
	Pairs *models.PairData
}

func (cs *CandlesCache) InsertCandle(candle *models.Candle) {

	for i := 1; i < 3; i++ {
		cs.Pairs.Five[i] = cs.Pairs.Five[i-1]
	}
	cs.Pairs.Five[0] = candle
}

func BuildCandleFromCandleList(candleList []*models.Candle) *models.Candle {
	tempCandle := &models.Candle{Open: 0, Close: 0, High: 0, Low: 0}

	for _, candle := range candleList {

		if tempCandle.High < candle.High {
			tempCandle.High = candle.High
		}
		if tempCandle.Low > candle.Low {
			tempCandle.Low = candle.Low
		}

	}
	tempCandle.Open = candleList[0].Open
	tempCandle.Close = candleList[len(candleList)-1].Close
	return tempCandle
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

func initializePairData() *models.PairData {

	return &models.PairData{
		Five:    make([]*models.Candle, 3),
		Fifteen: make([]*models.Candle, 2),
		Thirty:  make([]*models.Candle, 2),
		Hour:    make([]*models.Candle, 2)}
}

func New() *CandlesCache {

	return &CandlesCache{Pairs: initializePairData()}

}
