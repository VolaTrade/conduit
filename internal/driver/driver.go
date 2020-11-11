package driver

import (
	"github.com/google/wire"
	"github.com/volatrade/tickers/internal/service"
)

var Module = wire.NewSet(
	New,
)

type (
	TickersDriver struct {
		svc service.Service
	}
)

type Driver interface {
	RunListenerRoutines()
	Run()
	InitService()
}

func New(svc service.Service) *TickersDriver {
	return &TickersDriver{svc: svc}
}

//InitService initializes pairUrl list in cache and builds transactionChannels
func (td *TickersDriver) InitService() {
	if err := td.svc.BuildPairUrls(); err != nil {
		panic(err)
	}
	td.svc.BuildTransactionChannels(40)

}

func (td *TickersDriver) RunListenerRoutines() {

	for i := 0; i < 40; i++ {
		channel := td.svc.GetChannel(i)              //Gets channel for index
		go td.svc.ChannelListenAndHandle(channel, i) //Tells channels to listen for transaction data from sockets
	}
}

func (td *TickersDriver) Run() {

	go td.svc.CheckForDatabasePriveleges()
	sockets := td.svc.SpawnSocketRoutines(40)
	go td.svc.ReportRunning()
	for _, active_socket := range sockets {
		go td.svc.ConsumeTransferMessage(active_socket)
	}

	for {

	}

}
