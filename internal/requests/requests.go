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

type ConduitRequests struct {
	kstats stats.Stats
}

func New(stats stats.Stats) *ConduitRequests {
	return &ConduitRequests{kstats: stats}
}
