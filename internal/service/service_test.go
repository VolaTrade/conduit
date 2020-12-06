package service_test

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/volatrade/conduit/internal/cache"
	"github.com/volatrade/conduit/internal/mocks"
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/service"
	"github.com/volatrade/conduit/internal/stats"
)

type testSuite struct {
	mockController  *gomock.Controller
	mockConnections *mocks.MockConnections
	service         *service.ConduitService
	cache           cache.Cache
	client          *mocks.MockClient
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

//TODO remove return after stats updates 
func TestTransactionChannelsToCache(t *testing.T) {
	return 
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

//TODO remove return after stats updates 
func TestOrderBookChannelsToCache(t *testing.T) {
	return 
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

func TestCheckForDatabasePriveleges(t *testing.T) {
	ts := createTestSuite(t)

	ts.cache.InsertOrderBookRow(&models.OrderBookRow{
		Id:        123,
		Bids:      []byte("bids"),
		Asks:      []byte("asks"),
		Timestamp: time.Now(),
		Pair:      "BTCUSDT",
	})

	ts.cache.InsertTransaction(&models.Transaction{
		Id:        234,
		Pair:      "BTCUSDT",
		Price:     123.23,
		IsMaker:   false,
		Timestamp: time.Now(),
		Quantity:  12.21,
	})

	f, _ := os.Create("start")
	fmt.Printf("Created file %s", f.Name())
	ts.mockConnections.EXPECT().MakeConnections()
	ts.mockConnections.EXPECT().TransferTransactionCache(gomock.Any()).Return(nil).Times(1)
	ts.mockConnections.EXPECT().TransferOrderBookCache(gomock.Any()).Return(nil).Times(1)
	dir, _ := os.Getwd()

	println("Directory ---> ", dir)
	var wg sync.WaitGroup
	wg.Add(1)
	go ts.service.CheckForDatabasePriveleges(&wg)
	wg.Wait()

	assert.Equal(t, 0, ts.cache.TransactionsLength())
	assert.Equal(t, 0, ts.cache.TransactionsLength())
	os.Remove("start")
}

func TestBuildPairUrls(t *testing.T) {
	ts := createTestSuite(t)

	res := ts.service.BuildPairUrls()
	assert.True(t, res == nil)
}

func TestBuildOrderBookChannels(t *testing.T) {

	ts := createTestSuite(t)
	ts.service.BuildOrderBookChannels(3)
	c := ts.service.GetOrderBookChannel(2)
	assert.True(t, c != nil)
}

func TestBuildTransactionChannels(t *testing.T) {
	ts := createTestSuite(t)
	ts.service.BuildTransactionChannels(3)
	c := ts.service.GetTransactionChannel(2)
	assert.True(t, c != nil)
}

func TestSpawnSocketRoutines(t *testing.T) {
	ts := createTestSuite(t)

	sockets := ts.service.SpawnSocketRoutines(3)
	assert.True(t, sockets != nil)
}

func TestGetTransactionChannel(t *testing.T) {
	ts := createTestSuite(t)

	ts.service.BuildTransactionChannels(5)

	c := ts.service.GetTransactionChannel(4)

	assert.True(t, c != nil)
}

func TestGetOrderBookChannel(t *testing.T) {
	ts := createTestSuite(t)

	ts.service.BuildOrderBookChannels(5)

	c := ts.service.GetOrderBookChannel(4)

	assert.True(t, c != nil)
}
