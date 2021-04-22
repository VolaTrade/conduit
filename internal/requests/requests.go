//go:generate mockgen -package=mocks -destination=../mocks/requests.go github.com/volatrade/conduit/internal/requests Requests

package requests

import (
	"github.com/google/wire"
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
	cfg *Config
}

func New(cfg *Config) *ConduitRequests {
	return &ConduitRequests{cfg: cfg}
}
