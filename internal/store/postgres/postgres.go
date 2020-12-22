package postgres

import (
	"fmt"
	"log"

	"github.com/google/wire"
	_ "github.com/jackc/pgx/stdlib" //driver
	"github.com/jmoiron/sqlx"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
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
		kstats *stats.Stats
		logger *logger.Logger
	}
)

func New(cfg *Config, kstats *stats.Stats, logger *logger.Logger) *DB {
	postgres := &DB{config: cfg, kstats: kstats}

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
