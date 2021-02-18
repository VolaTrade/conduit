package test_redis

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/volatrade/conduit/internal/models"
)

func generateOrderBooksSlice(length int) []models.OrderBookRow {
	obs := make([]models.OrderBookRow, length)
	for i := 0; i < length; i++ {

		obs[i] = models.OrderBookRow{
			Id:        i + 1,
			Bids:      nil,
			Asks:      nil,
			Timestamp: time.Now(),
			Pair:      "btcusdt",
		}
	}

	return obs
}

func TestRollingCache(t *testing.T) {
	ts := createTestSuite(t)
	defer ts.endFunc()

	obs := generateOrderBooksSlice(31)

	for _, ob := range obs {
		if err := ts.cache.InsertOrderBookRowToRedis(&ob); err != nil {
			assert.NoError(t, err)
		}
	}

	rowsInRedis, err := ts.cache.GetOrderBookRowsFromRedis("btcusdt")

	assert.NoError(t, err)

	for _, strRow := range rowsInRedis {

		println("String row --->", strRow)
		assert.False(t, strings.Contains(strRow, "\"last_update_id\":1,"))
	}

}
