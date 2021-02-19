package test_redis

import (
	"log"
	"os"
	"testing"

	redis "github.com/volatrade/a-redis"
	"github.com/volatrade/conduit/internal/cache"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
)

var (
	redisConfig = &redis.Config{
		Host: "localhost",
		Port: 6379,
		DB:   0,
		Env:  "",
	}
)

type testSuite struct {
	cache   *cache.ConduitCache
	endFunc func()
	redis   redis.Redis
}

func createTestSuite(t *testing.T) *testSuite {
	noopLogger := logger.NewNoop()
	noopStats, _, _ := stats.NewNoop()
	redisClient, callback, err := redis.New(noopLogger, redisConfig, noopStats)

	if err != nil {
		log.Printf("Error trying to connect to redis : %e", err)
		os.Exit(1)
	}

	testCache := cache.New(noopLogger, redisClient)

	return &testSuite{
		cache:   testCache,
		endFunc: callback,
		redis:   redisClient,
	}
}

func TestMain(m *testing.M) {

	retCode := m.Run()

	os.Exit(retCode)

}
