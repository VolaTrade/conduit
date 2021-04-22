// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/volatrade/a-redis"
	"github.com/volatrade/conduit/internal/cache"
	"github.com/volatrade/conduit/internal/config"
	"github.com/volatrade/conduit/internal/cortex"
	"github.com/volatrade/conduit/internal/requests"
	"github.com/volatrade/conduit/internal/session"
	"github.com/volatrade/conduit/internal/store"
	"github.com/volatrade/conduit/internal/streamprocessor"
	"github.com/volatrade/currie-logs"
	"github.com/volatrade/k-stats"
	"github.com/volatrade/utilities/slack"
)

// Injectors from wire.go:

func InitializeAndRun(cfg config.FilePath) (streamprocessor.StreamProcessor, func(), error) {
	configConfig := config.NewConfig(cfg)
	postgresConfig := config.NewDBConfig(configConfig)
	statsConfig := config.NewStatsConfig(configConfig)
	statsStats, cleanup, err := stats.New(statsConfig)
	if err != nil {
		return nil, nil, err
	}
	loggerConfig := config.NewLoggerConfig(configConfig)
	loggerLogger, cleanup2, err := logger.New(loggerConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	sessionConfig := config.NewSessionConfig(configConfig)
	conduitSession := session.New(loggerLogger, sessionConfig, statsStats)
	conduitStorageConnections, cleanup3 := store.New(postgresConfig, statsStats, loggerLogger, conduitSession)
	aredisConfig := config.NewRedisConfig(configConfig)
	redis, cleanup4, err := aredis.New(loggerLogger, aredisConfig, statsStats)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	conduitCache := cache.New(loggerLogger, redis)
	requestsConfig := config.NewRequestConfig(configConfig)
	conduitRequests := requests.New(statsStats, requestsConfig)
	slackConfig := config.NewSlackConfig(configConfig)
	slackLogger := slack.New(slackConfig)
	cortexConfig := config.NewCortexConfig(configConfig)
	cortexClient, cleanup5, err := cortex.New(cortexConfig, statsStats, loggerLogger)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	conduitStreamProcessor, cleanup6 := streamprocessor.New(conduitStorageConnections, conduitCache, conduitRequests, conduitSession, statsStats, slackLogger, loggerLogger, cortexClient)
	return conduitStreamProcessor, func() {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

//cortexModule binds Cortex interface with ConduitCortex struct from session package
var cortexModule = wire.NewSet(cortex.Module, wire.Bind(new(cortex.Cortex), new(*cortex.CortexClient)))

//sessionModule binds Session interface with ConduitSession struct from session package
var sessionModule = wire.NewSet(session.Module, wire.Bind(new(session.Session), new(*session.ConduitSession)))

//cacheModule binds Cache interface with ConduitCache struct from Cache package
var cacheModule = wire.NewSet(cache.Module, wire.Bind(new(cache.Cache), new(*cache.ConduitCache)))

//StreamModule binds StreamProcessor interface with ConduitStreamProcessor struct from StreamProcessor package
var streamModule = wire.NewSet(streamprocessor.Module, wire.Bind(new(streamprocessor.StreamProcessor), new(*streamprocessor.ConduitStreamProcessor)))

//storageModule binds StorageConnections interface with ConduitStorageConnections struct from Store package
var storageModule = wire.NewSet(store.Module, wire.Bind(new(store.StorageConnections), new(*store.ConduitStorageConnections)))

//requestsModule module binds Requests interface with ConduitRequests struct from requests package
var requestsModule = wire.NewSet(requests.Module, wire.Bind(new(requests.Requests), new(*requests.ConduitRequests)))

//slack module binds Slack interface with SlackLogger struct from github.com/volatrade/utilities package
var slackModule = wire.NewSet(slack.Module, wire.Bind(new(slack.Slack), new(*slack.SlackLogger)))
