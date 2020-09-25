package models

type (
	// Candle defines the basic components of a candle
	Candle struct {
		Open  float64 `json:"open"`
		Close float64 `json:"close"`
		High  float64 `json:"high"`
		Low   float64 `json:"low"`
	}

	// PairData defines how a pair is structured
	PairData struct {
		Five    []*Candle // 3
		Fifteen []*Candle // 2
		Thirty  []*Candle // 2
		Hour    []*Candle // 1
	}

	// DynamoCandleItem defines what makes up a candle within dynamodb
	DynamoCandleItem struct {
		Timestamp string // Pkey
		PairName  string //Skey
		PairData  *PairData
	}
)
