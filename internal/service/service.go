//go:generate mockgen -package=mocks -destination=../mocks/service.go github.com/volatrade/conduit/internal/service Service
package service

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/cache"
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/requests"
	"github.com/volatrade/conduit/internal/socket"
	"github.com/volatrade/conduit/internal/stats"
	"github.com/volatrade/conduit/internal/store"
	log "github.com/volatrade/currie-logs"
	"github.com/volatrade/utilities/slack"
)

var (
	Module = wire.NewSet(
		New,
	)
)

type (
	Service interface {
		BuildPairUrls() error
		BuildTransactionChannels(count int)
		BuildOrderBookChannels(count int)
		CheckForDatabasePriveleges(wg *sync.WaitGroup)
		CheckForExit(wg *sync.WaitGroup, exit func())
		ListenAndHandle(queue chan *models.Transaction, obQueue chan *models.OrderBookRow, index int, wg *sync.WaitGroup, ch chan bool)
		SpawnSocketRoutines(psqlCount int) []*socket.BinanceSocket
		GetSocketsArrayLength() int
		GetTransactionChannel(index int) chan *models.Transaction
		GetOrderBookChannel(index int) chan *models.OrderBookRow
		ReportRunning(wg *sync.WaitGroup, ctx context.Context)
	}

	ConduitService struct {
		logger              *log.Logger
		id                  string
		cache               cache.Cache
		dbStreams           store.StorageConnections
		requests            requests.Requests
		slack               slack.Slack
		statsd              *stats.StatsD
		transactionChannels []chan *models.Transaction
		orderBookChannels   []chan *models.OrderBookRow
		writeToDB           bool
	}
)

func New(conns store.StorageConnections, ch cache.Cache, requests requests.Requests, stats *stats.StatsD, slackz slack.Slack, logger *log.Logger) *ConduitService {

	id := fmt.Sprintf("%d_%d", time.Now().Hour(), time.Now().Minute())

	logger.SetConstantField("Instance ID", id)
	return &ConduitService{
		cache:     ch,
		dbStreams: conns,
		requests:  requests,
		statsd:    stats,
		writeToDB: false,
		id:        id,
		slack:     slackz,
		logger:    logger,
	}
}

func (ts *ConduitService) ReportRunning(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	ts.statsd.Client.Gauge(fmt.Sprintf("conduit.instances.%s", ts.id), 1)

	for {

		select {

		case <-ctx.Done():
			ts.statsd.Client.Gauge(fmt.Sprintf("conduit.instances.%s", ts.id), 0)
			return

		}
	}
}

//TODO there's a better way to structure this
func (ts *ConduitService) CheckForDatabasePriveleges(wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	for {
		if _, writeToCache := os.Stat("start"); writeToCache == nil {
			ts.logger.Infow("establishing database connections, moving cache to databse, and purging cache")
			ts.dbStreams.MakeConnections()
			ts.writeToDB = true
			cachedTransactions := ts.cache.GetAllTransactions()
			cachedOrderBookRows := ts.cache.GetAllOrderBookRows()
			err = ts.dbStreams.TransferTransactionCache(cachedTransactions)
			if err != nil {
				panic(err)
			}

			err = ts.dbStreams.TransferOrderBookCache(cachedOrderBookRows)

			if err != nil {
				panic(err)
			}

			ts.cache.Purge()
			return
		}
	}
}

func (ts *ConduitService) CheckForExit(wg *sync.WaitGroup, exit func()) {
	defer wg.Done()
	for {
		if _, err := os.Stat("finish"); err == nil {
			exit()
			return
		}
	}
}

//Init reads all trading pairs from Binance and then proceeds to store them as keys in cache
func (ts *ConduitService) BuildPairUrls() error {

	tradingPairs, err := ts.requests.GetActiveBinanceExchangePairs()

	if err != nil {
		return err
	}

	for _, pair := range tradingPairs {

		if pair == "btcusdt" || pair == "ethusdt" || pair == "xrpusdt" {
			ts.cache.InsertPair(pair)
		}
	}

	return nil
}

//BuildTransactionChannels makes a slice of transaction struct channels
func (ts *ConduitService) BuildTransactionChannels(size int) {
	queues := make([]chan *models.Transaction, size)
	for i := 0; i < size; i++ {
		queue := make(chan *models.Transaction, 0)
		queues[i] = queue
	}
	ts.transactionChannels = queues
}

//BuildOrderBookChannels makes a slice of orderbook struct channels
func (ts *ConduitService) BuildOrderBookChannels(size int) {
	queues := make([]chan *models.OrderBookRow, size)

	for i := 0; i < size; i++ {
		queue := make(chan *models.OrderBookRow, 0)
		queues[i] = queue
	}

	ts.orderBookChannels = queues
}

func (ts *ConduitService) GetUrlsAndPair(index int) (string, string, string) {
	transactionURL, orderBookURL, err := ts.cache.GetTransactionOrderBookUrls(index)

	if err != nil {
		panic(err)
	}

	pair, err := ts.cache.GetPair(index)

	if err != nil {
		panic(err)
	}

	return transactionURL, orderBookURL, pair
}

func (ts *ConduitService) SpawnSocketRoutines(psqlCount int) []*socket.BinanceSocket {

	sockets := make([]*socket.BinanceSocket, 0)
	j := 0
	for i := 0; i < ts.cache.PairsLength(); i++ {
		if j >= psqlCount {
			j = 0
		}

		transactionURL, orderBookURL, pair := ts.GetUrlsAndPair(i)

		temp_stats := stats.StatsD{}                 //fix me
		temp_stats.Client = ts.statsd.Client.Clone() // fix me.. I am uneccesary
		socket, err := socket.NewSocket(transactionURL, orderBookURL, pair, ts.transactionChannels[j], ts.orderBookChannels[j], &temp_stats, ts.logger)

		if err != nil {
			fmt.Println(err)

		}
		sockets = append(sockets, socket)
		j++
	}

	return sockets
}

func (ts *ConduitService) GetSocketsArrayLength() int {
	return ts.cache.PairsLength()
}

func (ts *ConduitService) handleTransaction(tx *models.Transaction, index int) {
	if ts.writeToDB {
		ts.dbStreams.InsertTransactionToDataBase(tx, index)
		ts.statsd.Client.Increment(".conduit.sqlinserts.tx")

	} else {
		ts.cache.InsertTransaction(tx)
		ts.statsd.Client.Increment(".conduit.cacheinserts.tx")

	}
}

func (ts *ConduitService) handleOrderBookRow(tx *models.OrderBookRow, index int) {
	if ts.writeToDB {
		ts.dbStreams.InsertOrderBookRowToDataBase(tx, index)
		ts.statsd.Client.Increment(".conduit.sqlinserts.ob")

	} else {
		ts.cache.InsertOrderBookRow(tx)
		ts.statsd.Client.Increment(".conduit.cacheinserts.ob")

	}
}

func (ts *ConduitService) ListenAndHandle(txChannel chan *models.Transaction, obChannel chan *models.OrderBookRow, index int, wg *sync.WaitGroup, quit chan bool) {
	defer wg.Done()
	for {
		select {

		case <-quit:
			return

		case transaction := <-txChannel:
			ts.handleTransaction(transaction, index)

		case orderBookRow := <-obChannel:
			ts.handleOrderBookRow(orderBookRow, index)
		}

	}
}

func (ts *ConduitService) GetTransactionChannel(index int) chan *models.Transaction {
	return ts.transactionChannels[index]
}

func (ts *ConduitService) GetOrderBookChannel(index int) chan *models.OrderBookRow {
	return ts.orderBookChannels[index]
}

//TODO go routine grafana metric
