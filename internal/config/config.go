package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/volatrade/candles/internal/storage/postgres"
)

type Config struct {
	DbConfig postgres.Config
}

type FilePath string

func NewConfig(fileName FilePath) *Config {

	if err := godotenv.Load(string(fileName)); err != nil {
		log.Printf("Config file not found")
		log.Fatal(err)
	}

	return &Config{
		DbConfig: postgres.Config{
			Host:     os.Getenv("HOST"),
			Port:     os.Getenv("PORT"),
			Database: os.Getenv("DATABASE"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("PASSWORD"),
		},
	}
}

func NewDBConfig(cfg *Config) *postgres.Config {
	log.Println("Database config ---> ", cfg.DbConfig)
	return &cfg.DbConfig

}
