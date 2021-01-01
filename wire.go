//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/cache"
	"github.com/volatrade/conduit/internal/config"
	"github.com/volatrade/conduit/internal/requests"
	"github.com/volatrade/conduit/internal/session"

	"github.com/volatrade/conduit/internal/cortex"
	"github.com/volatrade/conduit/internal/store"
	sp "github.com/volatrade/conduit/internal/streamprocessor"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
	"github.com/volatrade/utilities/slack"
)

//sessionModule binds Session interface with ConduitSession struct from session package
var sessionModule = wire.NewSet(
	session.Module,
	wire.Bind(new(session.Session), new(*session.ConduitSession)),
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

//storageModule binds StorageConnections interface with ConduitStorageConnections struct from Store package
var cortexModule = wire.NewSet(
	cortex.Module,
	wire.Bind(new(cortex.Cortex), new(*cortex.CortexClient)),
)

//requestsModule module binds Requests interface with ConduitRequests struct from requests package
var requestsModule = wire.NewSet(
	requests.Module,
	wire.Bind(new(requests.Requests), new(*requests.ConduitRequests)),
)

//slack module binds Slack interface with SlackLogger struct from github.com/volatrade/utilities package
var slackModule = wire.NewSet(
	slack.Module,
	wire.Bind(new(slack.Slack), new(*slack.SlackLogger)),
)

func InitializeAndRun(cfg config.FilePath) (sp.StreamProcessor, func(), error) {

	panic(
		wire.Build(
			config.NewConfig,
			config.NewSessionConfig,
			config.NewDBConfig,
			config.NewStatsConfig,
			config.NewSlackConfig,
			config.NewLoggerConfig,
			logger.New,
			stats.New,
			cortexModule,
			sessionModule,
			storageModule,
			slackModule,
			requestsModule,
			cacheModule,
			streamModule,
		),
	)
}
