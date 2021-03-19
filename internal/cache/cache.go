//go:generate mockgen -package=mocks -destination=../mocks/cache.go github.com/volatrade/conduit/internal/cache Cache

package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/google/wire"
	"github.com/kr/pretty"
	redis "github.com/volatrade/a-redis"
	"github.com/volatrade/conduit/internal/models"
	log "github.com/volatrade/currie-logs"
)

var Module = wire.NewSet(
	New,
)

const (
	BASE_SOCKET_URL     string = "stream.binance.com:9443"
	OUT_OF_BOUNDS_ERROR string = "Index does not exist for pair slice"
)

type (
	Cache interface {
		GetAllOrderBookRows() []*models.OrderBookRow
		InsertOrderBookRow(obRow *models.OrderBookRow)
		InsertEntry(pair string)
		GetEntries() []*models.CacheEntry
		OrderBookRowsLength() int
		PurgeOrderBookRows()
		RowValidForCortex(pair string) bool
		InsertOrderBookRowToRedis(ob *models.OrderBookRow) error
		GetOrderBookRowsFromRedis(key string) ([]string, error)
	}

	ConduitCache struct {
		aredis        redis.Redis
		cortexObPairs *models.OrderBookPairs
		logger        *log.Logger
		entries       []*models.CacheEntry
		orderBookData []*models.OrderBookRow
		obMux         sync.Mutex
	}
)

//New ... constructor
func New(logger *log.Logger, aredis redis.Redis) *ConduitCache {

	return &ConduitCache{
		aredis:        aredis,
		logger:        logger,
		entries:       make([]*models.CacheEntry, 0),
		orderBookData: make([]*models.OrderBookRow, 0),
		cortexObPairs: &models.OrderBookPairs{Map: make(map[string]bool, 0)},
	}

}

//getOrderBookUrlString builds orderbook websocket url from pair
func getOrderBookUrlString(pair string) string {
	innerPath := fmt.Sprintf("ws/" + strings.ToLower(pair) + "@depth10@1000ms")
	socketUrl := url.URL{Scheme: "wss", Host: BASE_SOCKET_URL, Path: innerPath}
	return socketUrl.String()
}

//GetAllOrderBookRows returns cache slice of OrderBookRow model struct
func (cc *ConduitCache) GetAllOrderBookRows() []*models.OrderBookRow {
	return cc.orderBookData
}

//InsertEntry takes pair, builds URLs, appends data to Entry model struct, then adds struct to cache
func (cc *ConduitCache) InsertEntry(pair string) {

	entry := &models.CacheEntry{Pair: pair, ObUrl: getOrderBookUrlString(pair)}
	cc.entries = append(cc.entries, entry)

	if pair == "btcusdt" {
		cc.cortexObPairs.Map[pair] = true
	} else {
		cc.cortexObPairs.Map[pair] = false
	}
	println("POST")
	pretty.Print(cc.cortexObPairs)
}

//InsertOrderBookRow inserts OrderBookRow model struct to cache
func (cc *ConduitCache) InsertOrderBookRow(obRow *models.OrderBookRow) {
	if obRow == nil {
		cc.logger.Infow("Nil value passed in")
		return
	}

	cc.logger.Infow("cache insertion", "type", "orderbook snapshot", "cache length", cc.OrderBookRowsLength())
	cc.obMux.Lock()
	defer cc.obMux.Unlock()

	cc.orderBookData = append(cc.orderBookData, obRow)

}

func (cc *ConduitCache) RowValidForCortex(pair string) bool {

	
	println("Validating for", pair)
	if _, exists := cc.cortexObPairs.Map[pair]; !exists {
		cc.logger.Errorw(fmt.Sprintf("%s does not exist in OrderBookRows within memory cache", pair))
		return false
	}
	return cc.cortexObPairs.Map[pair]
}

func (cc *ConduitCache) PurgeOrderBookRows() {
	cc.orderBookData = nil

}

func (cc *ConduitCache) SetOrderBookPairs(obp *models.OrderBookPairs) {
	cc.cortexObPairs = obp
}

//GetEntries returns slice of CacheEntry struct
func (cc *ConduitCache) GetEntries() []*models.CacheEntry {
	return cc.entries
}

//OrderBookRowsLength used for testing && debuging
func (tc *ConduitCache) OrderBookRowsLength() int {

	if tc.orderBookData != nil {
		return len(tc.orderBookData)
	}
	return 0
}

func (cc *ConduitCache) InsertOrderBookRowToRedis(ob *models.OrderBookRow) error {

	println("Inserting into redis")

	err := cc.aredis.Ping(context.Background())

	if err != nil {
		println("Error pinging redis")
		println(err.Error())
	}
	bytes, err := json.Marshal(ob)

	if err != nil {
		return err
	}

	cc.logger.Infow("Redis insertion", "key", ob.Pair, "data", string(bytes))

	if err := cc.aredis.RPush(context.Background(), ob.Pair, bytes); err != nil {
		return err
	}
	println("Inserted")

	return nil
}

func (cc *ConduitCache) GetOrderBookRowsFromRedis(key string) ([]string, error) {
	obRows, err := cc.aredis.LRange(context.Background(), key, 0, -1)

	if err != nil {
		return nil, err
	}

	if len(obRows) < 30 {
		cc.logger.Infow("Redis orderbook list not long enough yet", "pair", key, "length", len(obRows))
		return nil, errors.New("List length in redis not long enough yet")
	}

	poppedVal, err := cc.aredis.LPop(context.Background(), key)

	if err != nil {
		return nil, err
	}

	cc.logger.Infow("Popped value from redis list", "value", poppedVal, "pair", key)
	return obRows[1:len(obRows)], nil

}
