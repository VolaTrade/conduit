package postgres

import (
	"fmt"
	"log"

	"github.com/google/wire"
	_ "github.com/jackc/pgx/stdlib" //driver
	"github.com/jmoiron/sqlx"
	"github.com/volatrade/candles/internal/models"
	"github.com/volatrade/candles/internal/stats"
)

const (
	INSERTION_QUERY = `INSERT INTO transactions(trade_id, time_stamp, pair, price, quantity, is_maker) VALUES(:id, :timestamp, :pair, :price, :quant, :maker) ON CONFLICT DO NOTHING`
)

var (
	Module = wire.NewSet(
		New,
	)
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

func New(cfg *Config, statsdClient *stats.StatsD) *DB {
	postgres := &DB{config: cfg, statz: statsdClient}

	return postgres
}

//Connect establishes connection to postgres server
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

//Close closes current connection w/ postgres server
func (postgres *DB) Close() error {
	if postgres != nil {
		err := postgres.DB.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//InsertTransaction inserts transaction into database
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

func (postgres *DB) BulkInsertCache(transactionList []*models.Transaction) error {

	tx := postgres.DB.MustBegin()
	stmt, err := tx.PrepareNamed(INSERTION_QUERY)
	if err != nil {
		return err
	}
	defer stmt.Close()

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

	return tx.Commit()

}
