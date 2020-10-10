package dynamo

import (
	"fmt"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

var Module = wire.NewSet(
	New,
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
	}
)

func New(config *Config) *DB {
	postgres := &DB{config: config}

	db, err := postgres.connect()
	if err != nil {
		panic(err)
	}
	postgres.DB = db
	return postgres
}

func (postgres *DB) connect() (*sqlx.DB, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		postgres.config.Host, postgres.config.Port, postgres.config.User, postgres.config.Password, postgres.config.Database)
	postgresDB, err := sqlx.Open("pgx", connString)
	if err != nil && postgresDB != nil {
		fmt.Println(err)
		return nil, err
	}
	err = postgresDB.Ping()
	if err != nil {
		print(fmt.Sprintf("postgres ping failed on startup, will keep trying. Error was %+v", err))
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
