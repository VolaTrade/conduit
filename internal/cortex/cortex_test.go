package cortex_test

import (
	"encoding/json"
	"time"

	"github.com/volatrade/conduit/internal/models"
)

var t = time.Now()
var socketMessage = []byte(`{"lastUpdateId":6688980481,"bids":[["17543.73000000","2.68375300"],["17543.52000000","0.00610000"],["17543.30000000","2.01579800"],["17542.86000000","0.01657800"],["17541.58000000","0.00610000"],["17541.14000000","0.79553600"],["17541.13000000","3.25000000"],["17541.12000000","0.01564000"],["17541.11000000","2.00000000"],["17541.00000000","0.03750000"]],"asks":[["17543.74000000","0.49208500"],["17543.75000000","0.00258200"],["17543.89000000","0.05696000"],["17543.90000000","0.00065900"],["17544.00000000","0.00100000"],["17544.01000000","0.00088000"],["17544.26000000","0.10000000"],["17544.32000000","0.00129200"],["17544.59000000","0.00064800"],["17544.67000000","0.00511200"]]}`)
var idealBids, _ = json.Marshal([][]string{{"17543.73000000", "2.68375300"}, {"17543.52000000", "0.00610000"}, {"17543.30000000", "2.01579800"}, {"17542.86000000", "0.01657800"}, {"17541.58000000", "0.00610000"}, {"17541.14000000", "0.79553600"}, {"17541.13000000", "3.25000000"}, {"17541.12000000", "0.01564000"}, {"17541.11000000", "2.00000000"}, {"17541.00000000", "0.03750000"}})
var idealAsks, _ = json.Marshal([][]string{{"17543.74000000", "0.49208500"}, {"17543.75000000", "0.00258200"}, {"17543.89000000", "0.05696000"}, {"17543.90000000", "0.00065900"}, {"17544.00000000", "0.00100000"}, {"17544.01000000", "0.00088000"}, {"17544.26000000", "0.10000000"}, {"17544.32000000", "0.00129200"}, {"17544.59000000", "0.00064800"}, {"17544.67000000", "0.00511200"}})

var idealOBRow = &models.OrderBookRow{
	Id:           6688980481,
	Bids:         idealBids,
	Asks:         idealAsks,
	Pair:         "BTCUSDT",
	CreationTime: t,
	Timestamp:    t.Format("2006:01:02 15:04"),
}

// type testSuite struct {x
// 	cfg          *cortex.Config
// 	cortexClient *cortex.CortexClient
// }

// func createTestSuite(t *testing.T) testSuite {
// 	mockLogger := logger.NewNoop()
// 	mockStats := kstats.NewNoop()
// 	cfg := &cortex.Config{Port: 9000, Host: "localhost"}
// 	cortexClient, _, _ := cortex.New(cfg, mockStats, mockLogger)
// 	return testSuite{
// 		cfg:          cfg,
// 		cortexClient: cortexClient,
// 	}
// }

// func TestMain(m *testing.M) {
// 	retCode := m.Run()
// 	os.Exit(retCode)
// }

// func TestSendOrderBookRow(t *testing.T) {
// 	ts := createTestSuite(t)
// 	err := ts.cortexClient.SendOrderBookRow(idealOBRow)
// 	assert.Nil(t, err)
// }
