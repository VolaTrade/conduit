//go:generate mockgen -package=mocks -destination=../mocks/store.go github.com/volatrade/conduit/internal/store StorageConnections

package store

import (
	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/store/postgres"
	stats "github.com/volatrade/k-stats"
)

var (
	Module = wire.NewSet(
		New,
	)
)

type (
	StorageConnections interface {
		MakeConnections()
		TransferTransactionCache(cacheData []*models.Transaction) error
		TransferOrderBookCache(cacheData []*models.OrderBookRow) error
		InsertTransactionToDataBase(transaction *models.Transaction, index int) error
		InsertOrderBookRowToDataBase(obRow *models.OrderBookRow, index int) error
	}

	ConduitStorageConnections struct {
		postgresConnections []*postgres.DB
	}
)

func New(cfg *postgres.Config, kstats *stats.Stats) *ConduitStorageConnections {
	arr := make([]*postgres.DB, 3)

	for i := 0; i < 3; i++ {
		tempDB := postgres.New(cfg, kstats)
		arr[i] = tempDB
	}

	return &ConduitStorageConnections{postgresConnections: arr}
}

func (ca *ConduitStorageConnections) MakeConnections() {

	for i := 0; i < 3; i++ {
		db, err := ca.postgresConnections[i].Connect()
		if err != nil {
			panic(err)
		}
		ca.postgresConnections[i].DB = db

	}

}

//TransferCache uses database connection at index 0 in connection array to transfer cache data to postgres
func (csc *ConduitStorageConnections) TransferTransactionCache(cacheData []*models.Transaction) error {
	return csc.postgresConnections[0].BulkInsertTransactions(cacheData)
}

func (csc *ConduitStorageConnections) TransferOrderBookCache(cacheData []*models.OrderBookRow) error {
	return csc.postgresConnections[0].BulkInsertOrderBookRows(cacheData)
}

func (csc *ConduitStorageConnections) InsertTransactionToDataBase(transaction *models.Transaction, index int) error {
	return csc.postgresConnections[index].InsertTransaction(transaction)
}

func (csc *ConduitStorageConnections) InsertOrderBookRowToDataBase(obRow *models.OrderBookRow, index int) error {

	return csc.postgresConnections[index].InsertOrderBookRow(obRow)
}
