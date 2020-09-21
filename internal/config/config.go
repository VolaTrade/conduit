package config

import (
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

	return &Config{}
}

func NewDBConfig(cfg *Config) *dynamo.Config {

	return &cfg.DbConfig

}
