//+build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/cache"
	"github.com/volatrade/conduit/internal/config"
	"github.com/volatrade/conduit/internal/conveyor"
	"github.com/volatrade/conduit/internal/requests"
	"github.com/volatrade/conduit/internal/session"

	"github.com/volatrade/conduit/internal/storage"
	sp "github.com/volatrade/conduit/internal/streamprocessor"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/go-grafana-graphite-client"
	"github.com/volatrade/utilities/slack"
)

//cacheModule binds Cache interface with ConduitCache struct from Cache package
var cacheModule = wire.NewSet(
	cache.Module,
	wire.Bind(new(cache.Cache), new(*cache.ConduitCache)),
)

//cacheModule binds Cache interface with ConduitCache struct from Cache package
var candleRequestModule = wire.NewSet(
	cache.Module,
	wire.Bind(new(cache.Cache), new(*cache.ConduitCache)),
)

//conveyorModule binds Conveyor interface with ConduitConveyor struct from conveyor package
var conveyorModule = wire.NewSet(
	conveyor.Module,
	wire.Bind(new(conveyor.Conveyor), new(*conveyor.ConduitConveyor)),
)

//requestsModule module binds Requests interface with ConduitRequests struct from requests package
var requestsModule = wire.NewSet(
	requests.Module,
	wire.Bind(new(requests.Requests), new(*requests.ConduitRequests)),
)

//sessionModule binds Session interface with ConduitSession struct from session package
var sessionModule = wire.NewSet(
	session.Module,
	wire.Bind(new(session.Session), new(*session.ConduitSession)),
)

//slack module binds Slack interface with SlackLogger struct from github.com/volatrade/utilities package
var slackModule = wire.NewSet(
	slack.Module,
	wire.Bind(new(slack.Slack), new(*slack.SlackLogger)),
)

//storageModule binds Store interface with ConduitStorage struct from Store package
var storageModule = wire.NewSet(
	storage.Module,
	wire.Bind(new(storage.Store), new(*storage.ConduitStorage)),
)

//StreamModule binds StreamProcessor interface with ConduitStreamProcessor struct from StreamProcessor package
var streamModule = wire.NewSet(
	sp.Module,
	wire.Bind(new(sp.StreamProcessor), new(*sp.ConduitStreamProcessor)),
)

func InitializeAndRun(ctx context.Context, cfg config.FilePath) (sp.StreamProcessor, func(), error) {

	panic(
		wire.Build(
			config.NewConfig,
			config.NewConveyorConfig,
			config.NewSessionConfig,
			config.NewDBConfig,
			config.NewStatsConfig,
			config.NewSlackConfig,
			config.NewLoggerConfig,
			config.NewLoggerOptions,
			config.NewRequestsConfig,
			logger.New,
			stats.NewClient,
			sessionModule,
			storageModule,
			slackModule,
			requestsModule,
			cacheModule,
			conveyorModule,
			streamModule,
		),
	)
}
