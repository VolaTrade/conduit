package driver

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/models"
	sproc "github.com/volatrade/conduit/internal/streamprocessor"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
)

var Module = wire.NewSet(
	New,
)

type (
	ConduitDriver struct {
		sp      sproc.StreamProcessor
		kstats  *stats.Stats
		session *models.Session
		logger  *logger.Logger
	}
)

type Driver interface {
	RunDataStreamListenerRoutines(ctx context.Context, wg *sync.WaitGroup)
	Run(ctx context.Context, wg *sync.WaitGroup, cancel func())
}

func New(sp sproc.StreamProcessor, stats *stats.Stats, sess *models.Session, logger *logger.Logger) *ConduitDriver {
	return &ConduitDriver{sp: sp, kstats: stats, session: sess, logger: logger}
}

func (td *ConduitDriver) RunDataStreamListenerRoutines(ctx context.Context, wg *sync.WaitGroup) {

	if err := td.sp.InsertPairsFromBinanceToCache(); err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go td.sp.ListenAndHandleDataChannels(ctx, i, wg)
	}
}

func (td *ConduitDriver) Run(ctx context.Context, wg *sync.WaitGroup, cancel func()) {
	wg.Add(3)
	go td.sp.ListenForDatabasePriveleges(ctx, wg)
	go td.sp.RunSocketRoutines(3)
	go td.sp.ListenForExit(ctx, wg, cancel)
	go td.reportRunning(ctx, wg)

}

func (cd *ConduitDriver) reportRunning(ctx context.Context, wg *sync.WaitGroup) {
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
