package streamprocessor_test

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	redis "github.com/volatrade/a-redis"
	"github.com/volatrade/conduit/internal/cache"
	cortex "github.com/volatrade/conduit/internal/cortex"
	"github.com/volatrade/conduit/internal/mocks"
	service "github.com/volatrade/conduit/internal/streamprocessor"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
)

type testSuite struct {
	mockController  *gomock.Controller
	mockConnections *mocks.MockStorageConnections
	service         *service.ConduitStreamProcessor
	cache           cache.Cache
	mockRequests    *mocks.MockRequests
	mockSession     *mocks.MockSession
	cortex          cortex.Cortex
}

func createTestSuite(t *testing.T) testSuite {
	mockController := gomock.NewController(t)
	redis, _, _ := redis.NewNoop()
	cache := cache.New(logger.NewNoop(), redis)

	stats, _, _ := stats.New(&stats.Config{Env: "DEV"})

	mockConnections := mocks.NewMockStorageConnections(mockController)

	mockRequests := mocks.NewMockRequests(mockController)
	mockSession := mocks.NewMockSession(mockController)

	cortexClient, _ := cortex.New(&cortex.Config{Port: 0}, stats, logger.NewNoop()) //This should be mocked

	mockSession.EXPECT().GetConnectionCount().Return(0).Times(100)
	svc, _ := service.New(mockConnections, cache, nil, mockSession, stats, nil, logger.NewNoop(), cortexClient)

	return testSuite{
		mockController:  mockController,
		service:         svc,
		mockConnections: mockConnections,
		cache:           cache,
		mockRequests:    mockRequests,
	}

}

func TestMain(m *testing.M) {

	retCode := m.Run()

	os.Exit(retCode)

}

// //TODO remove return after stats updates
// func TestTransactionChannelsToCache(t *testing.T) {
// 	ts := createTestSuite(t)

// 	ts.service.BuildTransactionChannels(1)
// 	ts.service.BuildOrderBookChannels(1)
// 	var wg sync.WaitGroup
// 	quit := make(chan bool)

// 	wg.Add(1)
// 	go ts.service.ListenAndHandleDataChannels(context.Background(), 0, &wg)
// 	txChannel := ts.service.GetTransactionChannel(0)

// 	for i := 0; i < 100; i++ {
// 		tx := &models.Transaction{}
// 		txChannel <- tx
// 	}
// 	quit <- true
// 	wg.Wait()
// 	assert.True(t, ts.cache.TransactionsLength() == 100)
// }

// //TODO remove return after stats updates
// func TestOrderBookChannelsToCache(t *testing.T) {
// 	ts := createTestSuite(t)

// 	ts.service.BuildTransactionChannels(1)
// 	ts.service.BuildOrderBookChannels(1)
// 	var wg sync.WaitGroup
// 	quit := make(chan bool)

// 	wg.Add(1)
// 	go ts.service.ListenAndHandleDataChannels(context.Background(), 1, &wg)
// 	obChannel := ts.service.GetOrderBookChannel(0)

// 	for i := 0; i < 100; i++ {
// 		ob := &models.OrderBookRow{}
// 		obChannel <- ob
// 	}
// 	quit <- true
// 	wg.Wait()
// 	assert.True(t, ts.cache.OrderBookRowsLength() == 100)

// }
