//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/volatrade/candles/internal/cache"
	"github.com/volatrade/candles/internal/config"
	"github.com/volatrade/candles/internal/driver"
	"github.com/volatrade/candles/internal/dynamo"
	"github.com/volatrade/candles/internal/service"
)

func InitializeAndRun(cfg config.FilePath) (*driver.CandlesDriver, error) {

	panic(
		wire.Build(
			config.NewConfig,
			config.NewDBConfig,
			storageModule,
			cacheModule,
			serviceModule,
			driver.New,
		),
	)
}

var cacheModule = wire.NewSet(
	cache.Module,
	wire.Bind(new(cache.Cache), new(*cache.CandlesCache)),
)

var serviceModule = wire.NewSet(
	service.Module,
	wire.Bind(new(service.Service), new(*service.CandlesService)),
)

var storageModule = wire.NewSet(
	dynamo.Module,
	wire.Bind(new(dynamo.Dynamo), new(*dynamo.DynamoSession)),
)
