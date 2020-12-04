package cache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volatrade/conduit/internal/cache"
	"github.com/volatrade/conduit/internal/models"
)

func TestInsertTransactionValueGetAndLength(t *testing.T) {
	c := cache.New()

	c.InsertTransaction(&models.Transaction{Price: 19})
	assert.True(t, c.TransactionsLength() == 1)
	assert.True(t, c.GetAllTransactions()[0].Price == 19)

	c.InsertTransaction(&models.Transaction{Price: 40})
	assert.True(t, 2 == c.TransactionsLength())
	assert.True(t, c.GetAllTransactions()[0].Price == 19)
	assert.True(t, c.GetAllTransactions()[1].Price == 40)

}

func TestInsertOrderBookValueGetAndLength(t *testing.T) {
	c := cache.New()

	c.InsertOrderBookRow(&models.OrderBookRow{Id: 19})
	assert.True(t, c.OrderBookRowsLength() == 1)
	assert.True(t, c.GetAllOrderBookRows()[0].Id == 19)

	c.InsertOrderBookRow(&models.OrderBookRow{Id: 34})
	assert.True(t, c.OrderBookRowsLength() == 2)
	assert.True(t, c.GetAllOrderBookRows()[0].Id == 19)
	assert.True(t, c.GetAllOrderBookRows()[1].Id == 34)

}

func TestPurge(t *testing.T) {
	c := cache.New()

	c.InsertOrderBookRow(&models.OrderBookRow{Id: 19})
	c.InsertOrderBookRow(&models.OrderBookRow{Id: 34})
	c.InsertTransaction(&models.Transaction{Price: 40})
	c.InsertTransaction(&models.Transaction{Price: 19})

	c.Purge()

	assert.True(t, c.GetAllOrderBookRows() == nil)
	assert.True(t, c.GetAllTransactions() == nil)
	assert.True(t, c.TransactionsLength() == 0)
	assert.True(t, c.OrderBookRowsLength() == 0)

}

func TestTransactionUrlsInsertAndGet(t *testing.T) {
	c := cache.New()
	pairs := []string{"ethusdt", "BTcUSdt"}

	for _, pair := range pairs {
		c.InsertPair(pair)
	}

	var txUrl string
	var odUrl string
	var err error

	txUrl, odUrl, err = c.GetTransactionOrderBookUrls(0)

	assert.Nil(t, err)
	assert.True(t, txUrl == "wss://stream.binance.com:9443/ws/ethusdt@trade")
	assert.True(t, odUrl == "wss://stream.binance.com:9443/ws/ethusdt@depth10@100ms")

	txUrl, odUrl, err = c.GetTransactionOrderBookUrls(1)

	assert.Nil(t, err)
	assert.True(t, txUrl == "wss://stream.binance.com:9443/ws/btcusdt@trade")
	assert.True(t, odUrl == "wss://stream.binance.com:9443/ws/btcusdt@depth10@100ms")

	txUrl, odUrl, err = c.GetTransactionOrderBookUrls(5)

	assert.True(t, txUrl == "")
	assert.True(t, odUrl == "")
	assert.True(t, err.Error() == cache.OUT_OF_BOUNDS_ERROR)

	assert.True(t, c.PairsLength() == 2)

}
