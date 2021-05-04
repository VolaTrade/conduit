//go:generate mockgen -package=mocks -destination=../mocks/requests.go github.com/volatrade/conduit/internal/requests Requests

package requests

import (
	"net/http"
	"time"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/models"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
)

var Module = wire.NewSet(
	New,
)

type Requests interface {
	GetActiveOrderbookPairs(retry int) ([]string, error)
	PostOrderbookRowToCortex(orderbookRow *models.OrderBookRow) error
}

type Config struct {
	BinanceApiUrl  string
	GatekeeperUrl  string
	RequestTimeout time.Duration
	CortexUrl      string
	CortexPort     int
}

type ConduitRequests struct {
	cfg    *Config
	statz  stats.Stats
	logger *logger.Logger
}

func New(cfg *Config, statz stats.Stats, logger *logger.Logger) *ConduitRequests {
	return &ConduitRequests{cfg: cfg, statz: statz, logger: logger}
}

func (cr *ConduitRequests) closeResponseBody(resp *http.Response) {

	if err := resp.Body.Close(); err != nil {
		cr.logger.Errorw("Error closing response")
	}
}
