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
	for i := 1; i < length; i++ {

		t := time.Now()
		obs[i] = models.OrderBookRow{
			Id:        i,
			Bids:      nil,
			Asks:      nil,
			Time:      t,
			Timestamp: t.Format("2006:01:02 15:04"),
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

		assert.False(t, strings.Contains(strRow, "\"last_update_id\":1,"))
	}

}
