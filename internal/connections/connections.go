package connections

import (
	"github.com/google/wire"
	"github.com/volatrade/tickers/internal/models"
	"github.com/volatrade/tickers/internal/postgres"
	"github.com/volatrade/tickers/internal/stats"
)

var Module = wire.NewSet(
	New,
)

type Connections interface {
	MakeConnections()
	TransferTransactionCache(cacheData []*models.Transaction) error
	TransferOrderBookCache(cacheData []*models.OrderBookRow) error
	InsertTransactionToDataBase(transaction *models.Transaction, index int) error
	InsertOrderBookRowToDataBase(obRow *models.OrderBookRow, index int) error
}

type ConnectionArray struct {
	Arr []*postgres.DB
}

func New(cfg *postgres.Config, statz *stats.StatsD) *ConnectionArray {
	arr := make([]*postgres.DB, 3)

	for i := 0; i < 3; i++ {
		temp_stats := stats.StatsD{}
		temp_stats.Client = statz.Client.Clone()
		tempDB := postgres.New(cfg, &temp_stats)
		arr[i] = tempDB
	}

	return &ConnectionArray{Arr: arr}
}

func (ca *ConnectionArray) MakeConnections() {

	for i := 0; i < 3; i++ {
		db, err := ca.Arr[i].Connect()
		if err != nil {
			panic(err)
		}
		ca.Arr[i].DB = db

	}

}

//TransferCache uses database connection at index 0 in connection array to transfer cache data to postgres
func (ca *ConnectionArray) TransferTransactionCache(cacheData []*models.Transaction) error {
	return ca.Arr[0].BulkInsertTransactions(cacheData)
}

func (ca *ConnectionArray) TransferOrderBookCache(cacheData []*models.OrderBookRow) error {
	return ca.Arr[0].BulkInsertOrderBookRows(cacheData)
}

//TODO add tarsnfer cache for OB

func (ca *ConnectionArray) InsertTransactionToDataBase(transaction *models.Transaction, index int) error {
	return ca.Arr[index].InsertTransaction(transaction)
}

func (ca *ConnectionArray) InsertOrderBookRowToDataBase(obRow *models.OrderBookRow, index int) error {
	println("Inserting order book data for index ->", index)
	println("list length -> ", len(ca.Arr))
	return ca.Arr[index].InsertOrderBookRow(obRow)
}
