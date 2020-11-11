package stats

import (
	"fmt"
	"runtime"
	"time"

	"gopkg.in/alexcesaro/statsd.v2"

	"github.com/google/wire"
)

var Module = wire.NewSet(
	New,
)

type (
	Stats interface {
		ReportGoRoutines()
		Increment(value string)
	}
	Config struct {
		Host string
		Port int
		Env  string
	}

	StatsD struct {
		Client *statsd.Client
	}
)

func New(cfg *Config) (*StatsD, error) {
	client, err := statsd.New(statsd.Address(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)), statsd.Prefix(cfg.Env))
	if err != nil {
		return nil, err
	}
	return &StatsD{Client: client}, nil
}

func ReportGoRoutines(statz *StatsD) {

	for {
		time.Sleep(1)
		statz.Client.Gauge("tickers.goroutines", runtime.NumGoroutine())
	}

}

//Encapsulation
func (statz *StatsD) Increment(value string) {
	statz.Client.Increment(value)
}
