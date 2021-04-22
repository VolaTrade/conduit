//go:generate mockgen -package=mocks -destination=../mocks/requests.go github.com/volatrade/conduit/internal/requests Requests

package requests

import (
	"github.com/google/wire"
	stats "github.com/volatrade/k-stats"
)

var Module = wire.NewSet(
	New,
)

type Requests interface {
	GetActiveOrderbookPairs() ([]string, error)
}

type Config struct {
	GatekeeperUrl string
}

type ConduitRequests struct {
	kstats stats.Stats
	cfg    *Config
}

func New(stats stats.Stats, cfg *Config) *ConduitRequests {
	return &ConduitRequests{kstats: stats, cfg: cfg}
}
