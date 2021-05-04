//go:generate mockgen -package=mocks -destination=../mocks/store.go github.com/volatrade/conduit/internal/store StorageConnections

package storage

import (
	"log"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/session"
	"github.com/volatrade/conduit/internal/storage/postgres"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
)

var (
	Module = wire.NewSet(
		New,
	)
)

type (
	Store interface {
		TransferOrderBookCache(cacheData []*models.OrderBookRow) error
	}

	ConduitStorage struct {
		kstats             stats.Stats
		session            session.Session
		postgresConnection *postgres.DB
	}
)

func New(cfg *postgres.Config, kstats stats.Stats, logger *logger.Logger, sess session.Session) (*ConduitStorage, func(), error) {
	pg := postgres.New(cfg, kstats, logger)

	conn, err := pg.Connect()

	if err != nil {
		return nil, nil, err
	}
	pg.DB = conn
	cs := &ConduitStorage{
		kstats:             kstats,
		session:            sess,
		postgresConnection: pg,
	}

	close := func() {
		log.Println("shutting down postgres connection")

		if err := pg.Close(); err != nil {
			log.Printf("Error obtained closing connection: %+v", err)
		}
		log.Println("postgres connection shutdown")
	}

	return cs, close, nil
}

func (cs *ConduitStorage) TransferOrderBookCache(cacheData []*models.OrderBookRow) error {

	if cacheData == nil {
		return nil
	}

	if err := cs.postgresConnection.BulkInsertOrderBookRows(cacheData); err != nil {
		return err
	}
	cs.kstats.Increment("postgres_transit", 1)
	return nil
}
