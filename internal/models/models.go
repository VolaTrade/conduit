package models

import (
	"fmt"
	"time"

	logger "github.com/volatrade/currie-logs"
)

type Session struct {
	ID string
}

func NewSession(logger *logger.Logger) *Session {
	id := fmt.Sprintf("%d_%d", time.Now().Hour(), time.Now().Minute())
	logger.SetConstantField("Session ID", id)
	return &Session{ID: id}

}


type CacheEntry struct  {
	TxUrl string 
	ObUrl string 
	Pair string 
}