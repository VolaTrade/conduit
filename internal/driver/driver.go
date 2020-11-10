package driver

import (
	"github.com/google/wire"
	"github.com/volatrade/candles/internal/service"
	"github.com/volatrade/candles/internal/socket"
)

var Module = wire.NewSet(
	New,
)

type (
	CandlesDriver struct {
		svc service.Service
	}
)

type Driver interface {
	Run()
}

func New(svc service.Service) *CandlesDriver {
	return &CandlesDriver{svc: svc}
}

func (cd *CandlesDriver) InitService() {
	if err := cd.svc.BuildPairUrls(); err != nil {
		panic(err)
	}
	cd.svc.BuildTransactionChannels(40)

}

func (cd *CandlesDriver) RunListenerRoutines() {

	for i := 0; i < 40; i++ {
		channel := cd.svc.GetChannel(i)
		go cd.svc.ChannelListenAndHandle(channel, i)
	}
}

func (cd *CandlesDriver) Run() {

	go cd.svc.CheckForDatabasePriveleges()
	sockets := cd.svc.SpawnSocketRoutines(40)

	for _, active_socket := range sockets {
		socket.ConsumeTransferMessage(active_socket)
	}

}
