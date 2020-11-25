package service_test

import (
	"os"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/volatrade/tickers/internal/cache"
	"github.com/volatrade/tickers/internal/mocks"
	"github.com/volatrade/tickers/internal/models"
	"github.com/volatrade/tickers/internal/service"
	"github.com/volatrade/tickers/internal/stats"
)

type testSuite struct {
	mockController  *gomock.Controller
	mockConnections *mocks.MockConnections
	service         *service.TickersService
	cache           cache.Cache
}

func createTestSuite(t *testing.T) testSuite {
	mockController := gomock.NewController(t)
	fakeStats, _ := stats.New(&stats.Config{Host: "localhost", Port: 8080, Env: "DEV"})
	cache := cache.New()
	mockConnections := mocks.NewMockConnections(mockController)

	svc := service.New(mockConnections, cache, nil, fakeStats, nil)

	return testSuite{
		mockController:  mockController,
		service:         svc,
		mockConnections: mockConnections,
		cache:           cache,
	}

}

func TestMain(m *testing.M) {

	retCode := m.Run()

	os.Exit(retCode)

}

func TestTransactionChannelsToCache(t *testing.T) {
	ts := createTestSuite(t)

	ts.service.BuildTransactionChannels(1)
	ts.service.BuildOrderBookChannels(1)
	var wg sync.WaitGroup
	quit := make(chan bool)

	wg.Add(1)
	go ts.service.ListenAndHandle(ts.service.GetTransactionChannel(0), ts.service.GetOrderBookChannel(0), 0, &wg, quit)
	println("HERE")
	txChannel := ts.service.GetTransactionChannel(0)

	for i := 0; i < 100; i++ {
		tx := &models.Transaction{}
		txChannel <- tx
	}
	quit <- true
	wg.Wait()
	assert.True(t, ts.cache.TransactionsLength() == 100)
}

func TestOrderBookChannelsToCache(t *testing.T) {
	ts := createTestSuite(t)

	ts.service.BuildTransactionChannels(1)
	ts.service.BuildOrderBookChannels(1)
	var wg sync.WaitGroup
	quit := make(chan bool)

	wg.Add(1)
	go ts.service.ListenAndHandle(ts.service.GetTransactionChannel(0), ts.service.GetOrderBookChannel(0), 0, &wg, quit)
	println("HERE")
	obChannel := ts.service.GetOrderBookChannel(0)

	for i := 0; i < 100; i++ {
		ob := &models.OrderBookRow{}
		obChannel <- ob
	}
	quit <- true
	wg.Wait()
	assert.True(t, ts.cache.OrderBookRowsLength() == 100)

}
