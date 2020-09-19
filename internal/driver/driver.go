package driver

import (
	"github.com/google/wire"
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

func (*CandlesDriver) Run() {

	for {
		println("Adrian hates vim")

	}

}
