//go:generate mockgen -package=mocks -destination=../mocks/requests.go github.com/volatrade/conduit/internal/requests Requests

package requests

import (
	"time"

	"github.com/google/wire"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
)

var Module = wire.NewSet(
	New,
)

type Requests interface {
	GetActiveOrderbookPairs(retry int) ([]string, error)
}

type Config struct {
	GatekeeperUrl  string
	RequestTimeout time.Duration
}

type ConduitRequests struct {
	cfg    *Config
	statz  *stats.Stats
	logger *logger.Logger
}

func New(cfg *Config, statz *stats.Stats, logger *logger.Logger) *ConduitRequests {
	return &ConduitRequests{cfg: cfg, statz: statz, logger: logger}
}
