package postgres

import (
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/stdlib" //driver

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

func (postgres *DB) InsertTransaction(mapping map[string]interface{}) error {

	query := `INSERT INTO transactions(time_stamp, pair, price, quantity, is_maker) VALUES($1, $2, $3, $4, $5);`

	log.Println(mapping["t"], mapping["T"], mapping["s"], mapping["p"], mapping["q"], mapping["m"])

	i := int64(mapping["T"].(float64)) / 1000
	tm := time.Unix(i, 0)
	res := postgres.DB.MustExec(query, tm, mapping["s"], mapping["p"], mapping["q"], mapping["m"])
	log.Println(res)
	return nil

}
