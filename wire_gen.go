// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

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
)

// Injectors from wire.go:

func InitializeAndRun(cfg config.FilePath) (driver.Driver, error) {
	configConfig := config.NewConfig(cfg)
	postgresConfig := config.NewDBConfig(configConfig)
	statsConfig := config.NewStatsConfig(configConfig)
	statsD, err := stats.New(statsConfig)
	if err != nil {
		return nil, err
	}
	connectionArray := connections.New(postgresConfig, statsD)
	tickersCache := cache.New()
	apiClient := client.New(statsD)
	tickersService := service.New(connectionArray, tickersCache, apiClient, statsD)
	tickersDriver := driver.New(tickersService, statsD)
	return tickersDriver, nil
}

// wire.go:

var cacheModule = wire.NewSet(cache.Module, wire.Bind(new(cache.Cache), new(*cache.TickersCache)))

var serviceModule = wire.NewSet(service.Module, wire.Bind(new(service.Service), new(*service.TickersService)))

var connectionModule = wire.NewSet(connections.Module, wire.Bind(new(connections.Connections), new(*connections.ConnectionArray)))

var apiClientModule = wire.NewSet(client.Module, wire.Bind(new(client.Client), new(*client.ApiClient)))

var driverModule = wire.NewSet(driver.Module, wire.Bind(new(driver.Driver), new(*driver.TickersDriver)))
