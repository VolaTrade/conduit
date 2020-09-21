package driver

import (
	"github.com/google/wire"
	"github.com/volatrade/candles/internal/cache"
	"github.com/volatrade/candles/internal/dynamo"
	"github.com/volatrade/candles/internal/service"
)

var Module = wire.NewSet(
	New,
)

type (
	CandlesDriver struct {
		svc *service.CandlesService
	}
)

type Driver interface {
	Run()
}

func New(service *service.CandlesService) *CandlesDriver {
	return &CandlesDriver{svc: service}
}

func (cd *CandlesDriver) Run() {

	if err := cd.svc.Init(); err != nil {
		panic(err)
	}
	//Insert concurrent workload distribution here

	// Test insertion
	candle, err := cache.NewCandle("1234", "1245", "1245", "12455", "timestamp2")
	if err != nil {
		panic(err)
	}

	dynamo, err := dynamo.New(&dynamo.Config{})
	if err != nil {
		panic(err)
	}

	err = dynamo.AddItem(candle, "candles")
	if err != nil {
		panic(err)
	}

}
