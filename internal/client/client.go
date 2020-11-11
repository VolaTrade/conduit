package client

import (
	"github.com/google/wire"
	"github.com/volatrade/tickers/internal/stats"
)

var Module = wire.NewSet(
	New,
)

type Client interface {
	GetActiveBinanceExchangePairs() ([]interface{}, error)
}

type ApiClient struct {
	statsd *stats.StatsD
}

func New(stats *stats.StatsD) *ApiClient {
	return &ApiClient{statsd: stats}
}
