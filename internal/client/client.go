package client

import (
	"os"
	"sync"

	"github.com/google/wire"
	"github.com/volatrade/candles/internal/models"
	"github.com/volatrade/candles/internal/stats"
	"github.com/volatrade/utilities/limiter"
)

var Module = wire.NewSet(
	New,
)

type Client interface {
	GetActiveBinanceExchangePairs() ([]interface{}, error)
	ConnectSocketAndReadTickData(socketUrl string, interrupt chan os.Signal, queue chan *models.Transaction, wg *sync.WaitGroup)
}

type ApiClient struct {
	rl     *limiter.RateLimiter
	statsd *stats.StatsD
}

func New(stats *stats.StatsD) *ApiClient {
	return &ApiClient{statsd: stats}
}
