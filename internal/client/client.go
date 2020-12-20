//go:generate mockgen -package=mocks -destination=../mocks/client.go github.com/volatrade/conduit/internal/client Client


package client

import (
	"github.com/google/wire"
	"github.com/volatrade/conduit/k-stats"
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
