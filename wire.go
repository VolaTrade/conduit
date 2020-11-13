//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/volatrade/tickers/internal/cache"
	"github.com/volatrade/tickers/internal/client"
	"github.com/volatrade/tickers/internal/config"
	"github.com/volatrade/tickers/internal/connections"
	"github.com/volatrade/tickers/internal/driver"
	"github.com/volatrade/tickers/internal/service"
	"github.com/volatrade/tickers/internal/stats"
	"github.com/volatrade/utilities/slack"
)

func InitializeAndRun(cfg config.FilePath) (driver.Driver, error) {

	panic(
		wire.Build(
			config.NewConfig,
			config.NewDBConfig,
			connectionModule,
			config.NewStatsConfig,
			config.NewSlackConfig,
			stats.New,
			slackModule,
			apiClientModule,
			cacheModule,
			serviceModule,
			driverModule,
		),
	)
}

var cacheModule = wire.NewSet(
	cache.Module,
	wire.Bind(new(cache.Cache), new(*cache.TickersCache)),
)

var serviceModule = wire.NewSet(
	service.Module,
	wire.Bind(new(service.Service), new(*service.TickersService)),
)

var connectionModule = wire.NewSet(
	connections.Module,
	wire.Bind(new(connections.Connections), new(*connections.ConnectionArray)),
)

var apiClientModule = wire.NewSet(
	client.Module,
	wire.Bind(new(client.Client), new(*client.ApiClient)),
)

var driverModule = wire.NewSet(
	driver.Module,
	wire.Bind(new(driver.Driver), new(*driver.TickersDriver)),
)

var slackModule = wire.NewSet(
	slack.Module,
	wire.Bind(new(slack.Slack), new(*slack.SlackLogger)),
)
