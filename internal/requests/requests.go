//go:generate mockgen -package=requests -destination=../mocks/requests.go github.com/volatrade/conduit/internal/requests Requests

package requests

import (
	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/stats"
)

var Module = wire.NewSet(
	New,
)

type Requests interface {
	GetActiveBinanceExchangePairs() ([]string, error)
}

type ConduitRequests struct {
	statsd *stats.StatsD
}

func New(stats *stats.StatsD) *ConduitRequests {
	return &ConduitRequests{statsd: stats}
}
