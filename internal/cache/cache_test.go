package cache_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volatrade/conduit/internal/cache"
	"github.com/volatrade/conduit/internal/models"
	log "github.com/volatrade/currie-logs"
)

var logger = log.NewNoop()

func TestInsertTransactionValueGetAndLength(t *testing.T) {
	c := cache.New(logger)

	c.InsertTransaction(&models.Transaction{Price: 19})
	assert.True(t, c.TransactionsLength() == 1)
	assert.True(t, c.GetAllTransactions()[0].Price == 19)

	c.InsertTransaction(&models.Transaction{Price: 40})
	assert.True(t, 2 == c.TransactionsLength())
	assert.True(t, c.GetAllTransactions()[0].Price == 19)
	assert.True(t, c.GetAllTransactions()[1].Price == 40)

}

func TestInsertOrderBookValueGetAndLength(t *testing.T) {
	c := cache.New(logger)

	c.InsertOrderBookRow(&models.OrderBookRow{Id: 19})
	assert.True(t, c.OrderBookRowsLength() == 1)
	assert.True(t, c.GetAllOrderBookRows()[0].Id == 19)

	c.InsertOrderBookRow(&models.OrderBookRow{Id: 34})
	assert.True(t, c.OrderBookRowsLength() == 2)
	assert.True(t, c.GetAllOrderBookRows()[0].Id == 19)
	assert.True(t, c.GetAllOrderBookRows()[1].Id == 34)

	c = cache.New(logger)

	for i := 1; i <= 39; i += 1 {
		c.InsertOrderBookRow(&models.OrderBookRow{Id: i})

	}
	assert.True(t, c.OrderBookRowsLength() == 39, fmt.Sprintf("asserting len(cache) =%d == 39", c.OrderBookRowsLength()))

}

func TestPurge(t *testing.T) {
	c := cache.New(logger)

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
