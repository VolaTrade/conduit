package session

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/wire"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
)

var (
	Module = wire.NewSet(
		New,
	)
)

type (
	Session interface {
		GetConnectionCount() int
		ReportRunning(ctx context.Context, wg *sync.WaitGroup)
	}

	Config struct {
		StorageConnections int
		Env                string
	}
	ConduitSession struct {
		id     string
		cfg    *Config
		kstats *stats.Stats
	}
)

func New(logger *logger.Logger, cfg *Config, kstats *stats.Stats) *ConduitSession {
	id := fmt.Sprintf("%d_%d", time.Now().Hour(), time.Now().Minute())
	logger.SetConstantField("ConduitSession ID", id)
	logger.SetConstantField("environment", cfg.Env)

	return &ConduitSession{id: id, cfg: cfg, kstats: kstats}

}

func (cs *ConduitSession) ReportRunning(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	cs.kstats.Gauge(fmt.Sprintf("conduit.instances.%s", cs.id), 1.0)

	for {

		select {

		case <-ctx.Done():
			println("Reporting zero")
			cs.kstats.Gauge(fmt.Sprintf("conduit.instances.%s", cs.id), 0.0)
			return

		}
	}
}

func (cs *ConduitSession) GetConnectionCount() int {
	return cs.cfg.StorageConnections
}
