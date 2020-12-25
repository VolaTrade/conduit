//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/cache"
	"github.com/volatrade/conduit/internal/config"
	"github.com/volatrade/conduit/internal/driver"
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/requests"
	"github.com/volatrade/conduit/internal/store"
	sp "github.com/volatrade/conduit/internal/streamprocessor"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
	"github.com/volatrade/utilities/slack"
)

//cacheModule binds Cache interface with ConduitCache struct from Cache package
var cacheModule = wire.NewSet(
	cache.Module,
	wire.Bind(new(cache.Cache), new(*cache.ConduitCache)),
)

//StreamModule binds StreamProcessor interface with ConduitStreamProcessor struct from StreamProcessor package
var streamModule = wire.NewSet(
	sp.Module,
	wire.Bind(new(sp.StreamProcessor), new(*sp.ConduitStreamProcessor)),
)

//storageModule binds StorageConnections interface with ConduitStorageConnections struct from Store package
var storageModule = wire.NewSet(
	store.Module,
	wire.Bind(new(store.StorageConnections), new(*store.ConduitStorageConnections)),
)

//requestsModule module binds Requests interface with ConduitRequests struct from requests package
var requestsModule = wire.NewSet(
	requests.Module,
	wire.Bind(new(requests.Requests), new(*requests.ConduitRequests)),
)

//driver module binds Driver interface with ConduitDriver struct from driver package
var driverModule = wire.NewSet(
	driver.Module,
	wire.Bind(new(driver.Driver), new(*driver.ConduitDriver)),
)

//slack module binds Slack interface with SlackLogger struct from github.com/volatrade/utilities package
var slackModule = wire.NewSet(
	slack.Module,
	wire.Bind(new(slack.Slack), new(*slack.SlackLogger)),
)

func InitializeAndRun(cfg config.FilePath) (driver.Driver, func(), error) {

	panic(
		wire.Build(
			config.NewConfig,
			//config.NewDriverConfig,
			config.NewDBConfig,
			config.NewStatsConfig,
			config.NewSlackConfig,
			config.NewLoggerConfig,
			logger.New,
			stats.New,
			models.NewSession,
			storageModule,
			slackModule,
			requestsModule,
			cacheModule,
			streamModule,
			driverModule,
		),
	)
}
