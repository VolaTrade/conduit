//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/volatrade/candles/internal/cache"
	"github.com/volatrade/candles/internal/client"
	"github.com/volatrade/candles/internal/config"
	"github.com/volatrade/candles/internal/driver"
	"github.com/volatrade/candles/internal/service"
	"github.com/volatrade/candles/internal/storage"
)

func InitializeAndRun(cfg config.FilePath) (*driver.CandlesDriver, error) {

	panic(
		wire.Build(
			config.NewConfig,
			config.NewDBConfig,
			storageModule,
			apiClientModule,
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
	storage.Module,
	wire.Bind(new(storage.Store), new(*storage.ConnectionArray)),
)

var apiClientModule = wire.NewSet(
	client.Module,
	wire.Bind(new(client.Client), new(*client.ApiClient)),
)
