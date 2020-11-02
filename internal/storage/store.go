package storage

import (
	"github.com/google/wire"
	"github.com/volatrade/candles/internal/stats"
	"github.com/volatrade/candles/internal/storage/postgres"
)

var Module = wire.NewSet(
	New,
)

type Store interface {
}

type ConnectionArray struct {
	Arr []*postgres.DB
}

func New(cfg *postgres.Config, statz *stats.StatsD) *ConnectionArray {
	arr := make([]*postgres.DB, 40)

	for i := 0; i < 40; i++ {
		temp_stats := stats.StatsD{}
		temp_stats.Client = statz.Client.Clone()
		tempDB := postgres.New(cfg, &temp_stats)
		arr[i] = tempDB
	}

	return &ConnectionArray{Arr: arr}
}

func (ca *ConnectionArray) MakeConnections() {

	for i := 0; i < 40; i++ {
		db, err := ca.Arr[i].Connect()
		if err != nil {
			panic(err)
		}
		ca.Arr[i].DB = db

	}

}
