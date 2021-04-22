package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	redis "github.com/volatrade/a-redis"
	"github.com/volatrade/conduit/internal/cortex"
	"github.com/volatrade/conduit/internal/session"
	"github.com/volatrade/conduit/internal/store/postgres"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
	"github.com/volatrade/utilities/slack"
)

type Config struct {
	DbConfig      postgres.Config
	RedisConfig   redis.Config
	StatsConfig   stats.Config
	SlackConfig   slack.Config
	SessionConfig session.Config
	CortexConfig  cortex.Config
}

type FilePath string

func NewConfig(fileName FilePath) *Config {

	if err := godotenv.Load(string(fileName)); err != nil {
		log.Printf("Config file not found")
		log.Fatal(err)
	}

	env := os.Getenv("ENV")

	if env != "DEV" && env != "PRD" && env != "INTEG" {
		log.Println("ENV ==>", env)
		log.Fatal("ENV var in config.env isn't set properly")
	}

	return &Config{
		DbConfig: postgres.Config{
			Host:     os.Getenv("PG_HOST"),
			Port:     os.Getenv("PG_PORT"),
			Database: os.Getenv("PG_DATABASE"),
			User:     os.Getenv("PG_USER"),
			Password: os.Getenv("PG_PASSWORD"),
		},
		StatsConfig: stats.Config{
			Host: os.Getenv("STATS_HOST"),
			Port: convertToInt(os.Getenv("STATS_PORT")),
			Env:  env,
		},
		SlackConfig: slack.Config{
			ApiKey:   os.Getenv("SLACK_API_KEY"),
			Location: "conduit",
			Env:      env,
		},
		SessionConfig: session.Config{
			StorageConnections: convertToInt(os.Getenv("STORAGE_CONNECTION_COUNT")),
			Env:                env,
		},

		RedisConfig: redis.Config{
			Host:     os.Getenv("REDIS_HOST"),
			Password: os.Getenv("REDIS_PASSWORD"),
			Port:     convertToInt(os.Getenv("REDIS_PORT")),
			DB:       convertToInt(os.Getenv("REDIS_DB")),
			Env:      env,
		},
		CortexConfig: cortex.Config{
			Host: os.Getenv("CORTEX_HOST"),
			Port: convertToInt(os.Getenv("CORTEX_PORT")),
		},
	}
}

func NewDBConfig(cfg *Config) *postgres.Config {
	return &cfg.DbConfig

}

func NewStatsConfig(cfg *Config) *stats.Config {
	return &cfg.StatsConfig
}

func NewSlackConfig(cfg *Config) *slack.Config {
	return &cfg.SlackConfig
}

func NewLoggerConfig(cfg *Config) *logger.Config {
	return nil
}
func NewRedisConfig(cfg *Config) *redis.Config {
	return &cfg.RedisConfig
}

func NewSessionConfig(cfg *Config) *session.Config {
	return &cfg.SessionConfig
}

func convertToInt(str string) int {
	intRep, err := strconv.Atoi(str)

	if err != nil {
		panic(err)
	}

	return intRep
}
