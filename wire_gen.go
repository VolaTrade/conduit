// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

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
	"github.com/volatrade/conduit/internal/streamprocessor"
	"github.com/volatrade/currie-logs"
	"github.com/volatrade/k-stats"
	"github.com/volatrade/utilities/slack"
)

// Injectors from wire.go:

func InitializeAndRun(ctx context.Context, cfg config.FilePath) (streamprocessor.StreamProcessor, func(), error) {
	configConfig := config.NewConfig(cfg)
	postgresConfig := config.NewDBConfig(configConfig)
	statsConfig := config.NewStatsConfig(configConfig)
	loggerConfig := config.NewLoggerConfig(configConfig)
	v := config.NewLoggerOptions(configConfig)
	loggerLogger, cleanup, err := logger.New(loggerConfig, v...)
	if err != nil {
		return nil, nil, err
	}
	statsStats, cleanup2, err := stats.New(statsConfig, loggerLogger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	sessionConfig := config.NewSessionConfig(configConfig)
	conduitSession := session.New(loggerLogger, sessionConfig, statsStats)
	conduitStorage, cleanup3, err := storage.New(postgresConfig, statsStats, loggerLogger, conduitSession)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	conduitCache := cache.New(loggerLogger)
	conveyorConfig := config.NewConveyorConfig(configConfig)
	conduitConveyor := conveyor.New(conveyorConfig, ctx, loggerLogger, conduitCache, conduitStorage)
	requestsConfig := config.NewRequestsConfig(configConfig)
	conduitRequests := requests.New(requestsConfig, statsStats, loggerLogger)
	slackConfig := config.NewSlackConfig(configConfig)
	slackLogger := slack.New(slackConfig)
	conduitStreamProcessor, cleanup4 := streamprocessor.New(ctx, conduitStorage, conduitCache, conduitConveyor, conduitRequests, conduitSession, statsStats, slackLogger, loggerLogger)
	return conduitStreamProcessor, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

//cacheModule binds Cache interface with ConduitCache struct from Cache package
var cacheModule = wire.NewSet(cache.Module, wire.Bind(new(cache.Cache), new(*cache.ConduitCache)))

//conveyorModule binds Conveyor interface with ConduitConveyor struct from conveyor package
var conveyorModule = wire.NewSet(conveyor.Module, wire.Bind(new(conveyor.Conveyor), new(*conveyor.ConduitConveyor)))

//requestsModule module binds Requests interface with ConduitRequests struct from requests package
var requestsModule = wire.NewSet(requests.Module, wire.Bind(new(requests.Requests), new(*requests.ConduitRequests)))

//sessionModule binds Session interface with ConduitSession struct from session package
var sessionModule = wire.NewSet(session.Module, wire.Bind(new(session.Session), new(*session.ConduitSession)))

//slack module binds Slack interface with SlackLogger struct from github.com/volatrade/utilities package
var slackModule = wire.NewSet(slack.Module, wire.Bind(new(slack.Slack), new(*slack.SlackLogger)))

//storageModule binds Store interface with ConduitStorage struct from Store package
var storageModule = wire.NewSet(storage.Module, wire.Bind(new(storage.Store), new(*storage.ConduitStorage)))

//StreamModule binds StreamProcessor interface with ConduitStreamProcessor struct from StreamProcessor package
var streamModule = wire.NewSet(streamprocessor.Module, wire.Bind(new(streamprocessor.StreamProcessor), new(*streamprocessor.ConduitStreamProcessor)))
