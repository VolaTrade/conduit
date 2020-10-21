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

	arr[0].Exec("CREATE TABLE transactions(
		trade_id UUID NOT NULL DEFAULT uuid_generate_v4 (),
		time_stamp TIMESTAMP,
		pair VARCHAR,
		price NUMERIC,
		quantity NUMERIC,
		is_maker boolean
	);")

	return &ConnectionArray{Arr: arr}
}
