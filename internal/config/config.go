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

	port, err := strconv.Atoi(os.Getenv("STATS_PORT"))

	if err != nil {
		println("#1")
		log.Fatal(err)
	}

	cortex_port, err := strconv.Atoi(os.Getenv("CORTEX_PORT"))

	if err != nil {
		println("#2")
		log.Fatal(err)
	}

	connCount, err := strconv.Atoi(os.Getenv("STORAGE_CONNECTION_COUNT"))

	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		println("#3")
		log.Fatal(err)
	}

	redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatal(err)
	}
	env := os.Getenv("ENV")

	if env != "DEV" && env != "PRD" && env != "INTEG" {
		log.Println("ENV ==>", env)
		log.Fatal("ENV var in config.env isn't set properly")
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
			Env:  env,
		},
		SlackConfig: slack.Config{
			ApiKey:   os.Getenv("SLACK_API_KEY"),
			Location: "conduit",
			Env:      env,
		},
		SessionConfig: session.Config{
			StorageConnections: connCount,
			Env:                env,
		},

		RedisConfig: redis.Config{
			Host:     os.Getenv("REDIS_HOST"),
			Password: os.Getenv("REDIS_PASSWORD"),
			Port:     redisPort,
			DB:       redisDb,
			Env:      env,
		},
		CortexConfig: cortex.Config{
			Host: os.Getenv("CORTEX_HOST"),
			Port: cortex_port,
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

func NewSlackConfig(cfg *Config) *slack.Config {
	log.Println("Slack config --->", cfg.SlackConfig)
	return &cfg.SlackConfig
}

func NewLoggerConfig(cfg *Config) *logger.Config {
	return nil
}
func NewRedisConfig(cfg *Config) *redis.Config {
	log.Println("Redis config --->", cfg.RedisConfig)
	return &cfg.RedisConfig
}

func NewSessionConfig(cfg *Config) *session.Config {
	log.Println("Session config --->", cfg.SessionConfig)
	return &cfg.SessionConfig
}

func NewCortexConfig(cfg *Config) *cortex.Config {
	return &cfg.CortexConfig
}
