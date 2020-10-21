package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/volatrade/candles/internal/stats"
	"github.com/volatrade/candles/internal/storage/postgres"
)

type Config struct {
	DbConfig    postgres.Config
	StatsConfig stats.Config
}

type FilePath string

func NewConfig(fileName FilePath) *Config {

	if err := godotenv.Load(string(fileName)); err != nil {
		log.Printf("Config file not found")
		log.Fatal(err)
	}

	port, err := strconv.Atoi(os.Getenv("STATS_PORT"))

	if err != nil {
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
		StatsConfig: stats.Config{
			Host: os.Getenv("STATS_HOST"),
			Port: port,
		},
	}
}

func NewDBConfig(cfg *Config) *postgres.Config {
	log.Println("Database config ---> ", cfg.DbConfig)
	return &cfg.DbConfig

}

func NewStatsConfig(cfg *Config) *stats.Config {
	log.Println("Stats config --->", cfg.StatsConfig)
	return &cfg.StatsConfig
}
