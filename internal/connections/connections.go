package connections

import (
	"github.com/google/wire"
	"github.com/volatrade/candles/internal/models"
	"github.com/volatrade/candles/internal/postgres"
	"github.com/volatrade/candles/internal/stats"
)

var Module = wire.NewSet(
	New,
)

type Connections interface {
	MakeConnections()
	TransferCache(cacheData []*models.Transaction) error
	InsertTransactionToDataBase(transaction *models.Transaction, index int) error
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

//TransferCache uses database connection at index 0 in connection array to transfer cache data to postgres
func (ca *ConnectionArray) TransferCache(cacheData []*models.Transaction) error {
	return ca.Arr[0].BulkInsertCache(cacheData)
}

func (ca *ConnectionArray) InsertTransactionToDataBase(transaction *models.Transaction, index int) error {
	return ca.Arr[index].InsertTransaction(transaction)
}
