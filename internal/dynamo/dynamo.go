package dynamo

import (
	"github.com/google/wire"
)

var Module = wire.NewSet(
	New,
)

type Dynamo interface {
	Test()
}

type (
	Config struct {
		ConnectionString string
	}

	CandlesDynamo struct {
		config *Config
	}
)

func New(cfg *Config) (*CandlesDynamo, error) {
	return &CandlesDynamo{config: cfg}, nil

}

func (*CandlesDynamo) Test() {
	println(":)")
}
