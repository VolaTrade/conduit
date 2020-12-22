package service_test

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/volatrade/conduit/internal/models"
)

func TestListenForDatabasePriveleges(t *testing.T) {
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

	var wg sync.WaitGroup
	wg.Add(1)
	go ts.service.ListenForDatabasePriveleges(&wg)
	wg.Wait()

	assert.Equal(t, 0, ts.cache.TransactionsLength())
	assert.Equal(t, 0, ts.cache.TransactionsLength())
	os.Remove("start")
}


