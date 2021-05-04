//go:generate mockgen -package=mocks -destination=../mocks/session.go github.com/volatrade/conduit/internal/session Session
package session

import (
	"context"
	"fmt"
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
		ReportRunning(ctx context.Context)
	}

	Config struct {
		Env string
	}
	ConduitSession struct {
		id     string
		cfg    *Config
		kstats stats.Stats
		logger *logger.Logger
	}
)

func New(logger *logger.Logger, cfg *Config, kstats stats.Stats) *ConduitSession {
	id := fmt.Sprintf("%d_%d", time.Now().Hour(), time.Now().Minute())
	logger.SetConstantField("ConduitSession ID", id)
	logger.SetConstantField("environment", cfg.Env)

	return &ConduitSession{id: id, cfg: cfg, kstats: kstats, logger: logger}

}

func (cs *ConduitSession) ReportRunning(ctx context.Context) {
	cs.kstats.Gauge(fmt.Sprintf("instances.%s", cs.id), 1.0)

	for range ctx.Done() {
		println("Reporting zero")
		cs.kstats.Gauge(fmt.Sprintf("instances.%s", cs.id), 0.0)

		return
	}
}
