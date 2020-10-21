package stats

import (
	"fmt"

	"gopkg.in/alexcesaro/statsd.v2"

	"github.com/google/wire"
)

var Module = wire.NewSet(
	New,
)

type (
	Stats interface {
	}
	Config struct {
		Host string
		Port int
	}

	StatsD struct {
		Client *statsd.Client
	}
)

func New(cfg *Config) (*StatsD, error) {
	client, err := statsd.New(statsd.Address(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)))

	if err != nil {
		return nil, err
	}
	return &StatsD{Client: client}, nil
}
