package cache_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	redis "github.com/volatrade/a-redis"
	"github.com/volatrade/conduit/internal/cache"
	"github.com/volatrade/conduit/internal/models"
	log "github.com/volatrade/currie-logs"
)

var logger = log.NewNoop()

func TestInsertOrderBookValueGetAndLength(t *testing.T) {
	redis, _, _ := redis.NewNoop()
	c := cache.New(logger, redis)

	c.InsertOrderBookRow(&models.OrderBookRow{Id: 19})
	assert.True(t, c.OrderBookRowsLength() == 1)
	assert.True(t, c.GetAllOrderBookRows()[0].Id == 19)

	c.InsertOrderBookRow(&models.OrderBookRow{Id: 34})
	assert.True(t, c.OrderBookRowsLength() == 2)
	assert.True(t, c.GetAllOrderBookRows()[0].Id == 19)
	assert.True(t, c.GetAllOrderBookRows()[1].Id == 34)

	c = cache.New(logger, redis)

	for i := 1; i <= 39; i += 1 {
		c.InsertOrderBookRow(&models.OrderBookRow{Id: i})

	}
	assert.True(t, c.OrderBookRowsLength() == 39, fmt.Sprintf("asserting len(cache) =%d == 39", c.OrderBookRowsLength()))

}

func TestPurge(t *testing.T) {
	redis, _, _ := redis.NewNoop()
	c := cache.New(logger, redis)

	c.InsertOrderBookRow(&models.OrderBookRow{Id: 19})
	c.InsertOrderBookRow(&models.OrderBookRow{Id: 34})

	c.PurgeOrderBookRows()

	assert.True(t, c.GetAllOrderBookRows() == nil)
	assert.True(t, c.OrderBookRowsLength() == 0)

}
