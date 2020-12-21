//go:generate mockgen -package=mocks -destination=../mocks/client.go github.com/volatrade/conduit/internal/client Client

package client

import (
	"github.com/google/wire"
	stats "github.com/volatrade/k-stats"
)

var Module = wire.NewSet(
	New,
)

type Client interface {
	GetActiveBinanceExchangePairs() ([]interface{}, error)
}

type ApiClient struct {
	kstats *stats.Stats
}

func New(stats *stats.Stats) *ApiClient {
	return &ApiClient{kstats: stats}
}
