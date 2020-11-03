package postgres

import (
	"fmt"
	"log"

	_ "github.com/jackc/pgx/stdlib" //driver
	"github.com/volatrade/candles/internal/cache"
	"github.com/volatrade/candles/internal/models"
	"github.com/volatrade/candles/internal/stats"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

var Module = wire.NewSet(
	New,
)

const (
	INSERTION_QUERY = `INSERT INTO transactions(trade_id, time_stamp, pair, price, quantity, is_maker) VALUES(:id, :timestamp, :pair, :price, :quant, :maker) ON CONFLICT DO NOTHING`
)

type (
	Config struct {
		Host     string
		Port     string
		Database string
		User     string
		Password string
	}

	DB struct {
		DB     *sqlx.DB
		config *Config
		statz  *stats.StatsD
	}
)

func New(config *Config, statsdClient *stats.StatsD) *DB {
	postgres := &DB{config: config, statz: statsdClient}

	return postgres
}

func (postgres *DB) Connect() (*sqlx.DB, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		postgres.config.Host, postgres.config.Port, postgres.config.User, postgres.config.Password, postgres.config.Database)
	log.Println("Connection string -->", connString)
	postgresDB, err := sqlx.Open("pgx", connString)
	if err != nil && postgresDB != nil {
		log.Println("Error connecting to database")
		log.Println(err)
		return nil, err
	}
	err = postgresDB.Ping()
	if err != nil {
		log.Println(fmt.Sprintf("postgres ping failed on startup, will keep trying. Error was %+v", err))
	}
	return postgresDB, nil
}

func (postgres *DB) Close() error {
	if postgres != nil {
		err := postgres.DB.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (postgres *DB) InsertTransaction(transaction *models.Transaction) error {
	stmt, err := postgres.DB.PrepareNamed(INSERTION_QUERY)

	if err != nil {
		return err
	}

	defer stmt.Close()
	result, err := stmt.Exec(transaction)

	if err != nil {
		return err
	}
	if rows, err := result.RowsAffected(); rows == 0 && err == nil {
		postgres.statz.Client.Increment("duplicate_inserts")
	}

	return err

}

func (postgres *DB) PurgeCache(cache *cache.TickersCache) error {

	tx := postgres.DB.MustBegin()
	stmt, err := tx.PrepareNamed(INSERTION_QUERY)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, transactionList := range cache.Pairs {

		for _, transaction := range transactionList {

			if transaction == nil {
				continue
			}
			result, err := stmt.Exec(transaction)

			if err != nil {
				tx.Rollback()
				return err
			}

			if rows, err := result.RowsAffected(); rows == 0 && err == nil {
				postgres.statz.Client.Increment("duplicate_inserts")
			}
			if err != nil {
				tx.Rollback()
				return err
			}

		}
	}
	println("Cache has finished purgery")
	return tx.Commit()

}
