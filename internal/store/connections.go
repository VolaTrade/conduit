//go:generate mockgen -package=mocks -destination=../mocks/store.go github.com/volatrade/conduit/internal/store StorageConnections

package store

import (
	"log"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/session"
	"github.com/volatrade/conduit/internal/store/postgres"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
)

var (
	Module = wire.NewSet(
		New,
	)
)

type (
	StorageConnections interface {
		MakeConnections() error
		TransferOrderBookCache(cacheData []*models.OrderBookRow) error
		InsertOrderBookRowToDataBase(obRow *models.OrderBookRow, index int) error
	}

	ConduitStorageConnections struct {
		session             session.Session
		postgresConnections []*postgres.DB
	}
)

func New(cfg *postgres.Config, kstats *stats.Stats, logger *logger.Logger, sess session.Session) (*ConduitStorageConnections, func()) {
	arr := make([]*postgres.DB, sess.GetConnectionCount())

	for i := 0; i < sess.GetConnectionCount(); i++ {
		tempDB := postgres.New(cfg, kstats, logger)
		arr[i] = tempDB
	}

	close := func() {
		log.Println("shutting down postgres connections")

		for _, db := range arr {

			if db.DB == nil {
				return
			}

			if err := db.Close(); err != nil {
				log.Printf("Error obtained closing connection: %+v", err)
			}
		}

		log.Println("postgres connections shutdown")

	}

	return &ConduitStorageConnections{postgresConnections: arr, session: sess}, close
}

func (ca *ConduitStorageConnections) MakeConnections() error {

	log.Println("MAKING DATABASE CONNECTIONS")
	for i := 0; i < ca.session.GetConnectionCount(); i++ {
		db, err := ca.postgresConnections[i].Connect()
		if err != nil {
			return err
		}
		ca.postgresConnections[i].DB = db

	}
	return nil
}


func (csc *ConduitStorageConnections) TransferOrderBookCache(cacheData []*models.OrderBookRow) error {

	if cacheData == nil {
		return nil
	}
	return csc.postgresConnections[0].BulkInsertOrderBookRows(cacheData)
}

func (csc *ConduitStorageConnections) InsertOrderBookRowToDataBase(obRow *models.OrderBookRow, index int) error {

	return csc.postgresConnections[index].InsertOrderBookRow(obRow)
}
