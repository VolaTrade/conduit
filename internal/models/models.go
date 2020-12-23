package models

import (
	"fmt"
	"time"

	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
)

type Session struct {
	ID string
}

func NewSession(logger *logger.Logger, cfg *stats.Config) *Session {
	id := fmt.Sprintf("%d_%d", time.Now().Hour(), time.Now().Minute())
	logger.SetConstantField("Session ID", id)

	logger.SetConstantField("environment", cfg.Env)
	return &Session{ID: id}

}

type CacheEntry struct {
	TxUrl string
	ObUrl string
	Pair  string
}
