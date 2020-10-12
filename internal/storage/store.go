package storage

import (
	"github.com/google/wire"
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

func New(cfg *postgres.Config) *ConnectionArray {
	arr := make([]*postgres.DB, 40)

	for i := 0; i < 40; i++ {
		tempDB := postgres.New(cfg)
		arr[i] = tempDB
	}

	return &ConnectionArray{Arr: arr}
}
