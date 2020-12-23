package driver

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/service"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
)

var Module = wire.NewSet(
	New,
)

type (
	ConduitDriver struct {
		svc     service.Service
		kstats  *stats.Stats
		session *models.Session
		logger  *logger.Logger
	}
)

type Driver interface {
	RunDataStreamListenerRoutines(wg *sync.WaitGroup, ch chan bool)
	Run(wg *sync.WaitGroup)
}

func New(svc service.Service, stats *stats.Stats, sess *models.Session, logger *logger.Logger) *ConduitDriver {
	return &ConduitDriver{svc: svc, kstats: stats, session: sess, logger: logger}
}

func (td *ConduitDriver) RunDataStreamListenerRoutines(wg *sync.WaitGroup, ch chan bool) {

	if err := td.svc.InsertPairsFromBinanceToCache(); err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go td.svc.ListenAndHandleDataChannels(i, wg, ch)
	}
}

func (td *ConduitDriver) Run(wg *sync.WaitGroup) {

	ctx, cancel := context.WithCancel(context.Background())

	go td.svc.ListenForDatabasePriveleges(wg)
	go td.svc.RunSocketRoutines(3)
	go td.svc.ListenForExit(wg, cancel)
	go td.reportRunning(wg, ctx)

}

func (cd *ConduitDriver) reportRunning(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	cd.kstats.Gauge(fmt.Sprintf("conduit.instances.%s", cd.session.ID), 1.0) //should this Gauge be 1?

	for {

		select {

		case <-ctx.Done():
			println("Reporting zero")
			cd.kstats.Gauge(fmt.Sprintf("conduit.instances.%s", cd.session.ID), 0.0)
			return

		}
	}
}
