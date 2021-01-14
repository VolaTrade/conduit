package models_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/volatrade/conduit/internal/models"
)

var ts = time.Now()
var socketMessage = []byte(`{"lastUpdateId":6688980481,"bids":[["17543.73000000","2.68375300"],["17543.52000000","0.00610000"],["17543.30000000","2.01579800"],["17542.86000000","0.01657800"],["17541.58000000","0.00610000"],["17541.14000000","0.79553600"],["17541.13000000","3.25000000"],["17541.12000000","0.01564000"],["17541.11000000","2.00000000"],["17541.00000000","0.03750000"]],"asks":[["17543.74000000","0.49208500"],["17543.75000000","0.00258200"],["17543.89000000","0.05696000"],["17543.90000000","0.00065900"],["17544.00000000","0.00100000"],["17544.01000000","0.00088000"],["17544.26000000","0.10000000"],["17544.32000000","0.00129200"],["17544.59000000","0.00064800"],["17544.67000000","0.00511200"]]}`)
var idealBidsSl, _ = [][]string{{"17543.73000000", "2.68375300"}, {"17543.52000000", "0.00610000"}, {"17543.30000000", "2.01579800"}, {"17542.86000000", "0.01657800"}, {"17541.58000000", "0.00610000"}, {"17541.14000000", "0.79553600"}, {"17541.13000000", "3.25000000"}, {"17541.12000000", "0.01564000"}, {"17541.11000000", "2.00000000"}, {"17541.00000000", "0.03750000"}})
var idealAsksSl, _ = [][]string{{"17543.74000000", "0.49208500"}, {"17543.75000000", "0.00258200"}, {"17543.89000000", "0.05696000"}, {"17543.90000000", "0.00065900"}, {"17544.00000000", "0.00100000"}, {"17544.01000000", "0.00088000"}, {"17544.26000000", "0.10000000"}, {"17544.32000000", "0.00129200"}, {"17544.59000000", "0.00064800"}, {"17544.67000000", "0.00511200"}})
var idealOBRow = &models.OrderBookRow{
	Id:   6688980481,
	Bids: idealBidsSl,
	Asks: idealAsksSl,
	Pair: "BTCUSDT",
}

var socketTransaction = []byte(`{"e":"trade","E":1605862294342,"s":"BTCUSDT","t":473476704,"p":"18251.11000000","q":"0.08256400","b":3662513230,"a":3662513203,"T":1605862294341,"m":false,"M":true}`)

func TestUnmarshalOrderBook(t *testing.T) {

	rec, err := models.UnmarshalOrderBookJSON(socketMessage, "BTCUSDT")
	if err != nil {
		panic(err)
	}

	rec.Timestamp = ts
	idealOBRow.Timestamp = ts
	fmt.Println(rec)
	fmt.Println(idealOBRow)

	assert.True(t, err == nil)
	assert.EqualValues(t, idealOBRow, rec, "Values should match")
}

func TestUnmarshalTransactionJSON(t *testing.T) {
	message := socketTransaction
	ts := int64(1605862294341 / 1000)
	exp := &models.Transaction{Id: 3662513203, Timestamp: time.Unix(ts, 0), Pair: "BTCUSDT", Price: 18251.11000000, Quantity: 0.08256400, IsMaker: false}
	ret, err := models.UnmarshalTransactionJSON(message)

	assert.NoError(t, err)
	assert.EqualValues(t, exp, ret)
}
