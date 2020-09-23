package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/volatrade/candles/internal/dynamo"
)

type Config struct {
	DbConfig dynamo.Config
}

type FilePath string

func NewConfig(fileName FilePath) *Config {

	if err := godotenv.Load(string(fileName)); err != nil {
		panic(err)
	}

	return &Config{
		DbConfig: dynamo.Config{TableName: os.Getenv("TABLE_NAME")},
	}
}

func NewDBConfig(cfg *Config) *dynamo.Config {

	return &cfg.DbConfig

}
