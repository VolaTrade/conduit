package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/volatrade/conduit/internal/conveyor"
	"github.com/volatrade/conduit/internal/requests"
	"github.com/volatrade/conduit/internal/session"
	"github.com/volatrade/conduit/internal/storage/postgres"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
	"github.com/volatrade/utilities/slack"
)

type Config struct {
	ConveyorConfig conveyor.Config
	DbConfig       postgres.Config
	StatsConfig    stats.Config
	SlackConfig    slack.Config
	SessionConfig  session.Config
	RequestsConfig requests.Config
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

		ConveyorConfig: conveyor.Config{
			ShiftInterval: convertToInt(os.Getenv("DISPATCH_INTERVAL")),
		},
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
			Env: env,
		},
		RequestsConfig: requests.Config{
			GatekeeperUrl:  os.Getenv("GATEKEEPER_URL"),
			RequestTimeout: time.Duration(convertToInt(os.Getenv("REQUEST_TIMEOUT"))) * time.Second,
			CortexUrl:      os.Getenv("CORTEX_URL"),
			CortexPort:     convertToInt(os.Getenv("CORTEX_PORT")),
		},
	}
}

func NewRequestsConfig(cfg *Config) *requests.Config {
	log.Println("Requests config ---> ", cfg.RequestsConfig)
	return &cfg.RequestsConfig
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
func NewConveyorConfig(cfg *Config) *conveyor.Config {
	return &cfg.ConveyorConfig
}

func NewLoggerOptions(cfg *Config) []logger.Option {
	return []logger.Option{func(l *logger.Logger) {}}
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
