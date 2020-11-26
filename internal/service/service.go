package service

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/wire"
	"github.com/volatrade/tickers/internal/cache"
	"github.com/volatrade/tickers/internal/client"
	"github.com/volatrade/tickers/internal/connections"
	"github.com/volatrade/tickers/internal/models"
	"github.com/volatrade/tickers/internal/socket"
	"github.com/volatrade/tickers/internal/stats"
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
		ListenAndHandle(queue chan *models.Transaction, obQueue chan *models.OrderBookRow, index int, wg *sync.WaitGroup, ch chan bool)
		SpawnSocketRoutines(psqlCount int) []*socket.BinanceSocket
		GetSocketsArrayLength() int
		GetTransactionChannel(index int) chan *models.Transaction
		GetOrderBookChannel(index int) chan *models.OrderBookRow
		ReportRunning(wg *sync.WaitGroup)
	}

	TickersService struct {
		id                  string
		cache               cache.Cache
		connections         connections.Connections
		exch                client.Client
		slack               slack.Slack
		statsd              *stats.StatsD
		transactionChannels []chan *models.Transaction
		orderBookChannels   []chan *models.OrderBookRow
		writeToDB           bool
	}
)

func New(conns connections.Connections, ch cache.Cache, cl *client.ApiClient, stats *stats.StatsD, slackz slack.Slack) *TickersService {

	return &TickersService{
		cache:       ch,
		connections: conns,
		exch:        cl,
		statsd:      stats,
		writeToDB:   false,
		id:          fmt.Sprintf("%d_%d", time.Now().Hour(), time.Now().Minute()),
		slack:       slackz,
	}
}

func (ts *TickersService) ReportRunning(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		time.Sleep(10000)
		ts.statsd.Client.Increment(fmt.Sprintf("tickers.instances.%s", ts.id))
	}
}

//TODO there's a better way to structure this
func (ts *TickersService) CheckForDatabasePriveleges(wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	for {
		if _, writeToCache := os.Stat("start"); writeToCache == nil {
			log.Println("making connections to DB NOW")
			ts.connections.MakeConnections()
			ts.writeToDB = true
			log.Println("Purging cache")
			cachedTransactions := ts.cache.GetAllTransactions()
			cachedOrderBookRows := ts.cache.GetAllOrderBookRows()
			err = ts.connections.TransferTransactionCache(cachedTransactions)
			if err != nil {
				panic(err)
			}

			err = ts.connections.TransferOrderBookCache(cachedOrderBookRows)

			if err != nil {
				panic(err)
			}

			ts.cache.Purge()
			//TODO insert transfer logic for order book data
			return
		}

	}
}

//Init reads all trading pairs from Binance and then proceeds to store them as keys in cache
func (ts *TickersService) BuildPairUrls() error {
	tradingCryptosList, err := ts.exch.GetActiveBinanceExchangePairs()
	if err != nil {
		return err
	}

	for _, val := range tradingCryptosList {
		temp := val.(map[string]interface{}) //type casting
		id := strings.ToLower(temp["symbol"].(string))

		if id == "btcusdt" || id == "ethusdt" || id == "xrpusdt" {
			ts.cache.InsertPair(id)
		}
	}

	return nil
}

//BuildTransactionChannels makes a slice of channels
func (ts *TickersService) BuildTransactionChannels(count int) {
	queues := make([]chan *models.Transaction, count)
	for i := 0; i < count; i++ {
		queue := make(chan *models.Transaction, 0)
		queues[i] = queue
	}
	ts.transactionChannels = queues
}

func (ts *TickersService) BuildOrderBookChannels(count int) {
	queues := make([]chan *models.OrderBookRow, count)

	for i := 0; i < count; i++ {
		queue := make(chan *models.OrderBookRow, 0)
		queues[i] = queue
	}

	ts.orderBookChannels = queues
}

func (ts *TickersService) GetUrlsAndPair(index int) (string, string, string) {
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

func (ts *TickersService) SpawnSocketRoutines(psqlCount int) []*socket.BinanceSocket {

	sockets := make([]*socket.BinanceSocket, 0)
	j := 0
	for i := 0; i < ts.cache.PairsLength(); i++ {
		if j >= psqlCount {
			j = 0
		}

		transactionURL, orderBookURL, pair := ts.GetUrlsAndPair(i)

		temp_stats := stats.StatsD{}                 //fix me
		temp_stats.Client = ts.statsd.Client.Clone() // fix me.. I am uneccesary
		socket, err := socket.NewSocket(transactionURL, orderBookURL, pair, ts.transactionChannels[j], ts.orderBookChannels[j], &temp_stats)

		if err != nil {
			fmt.Println(err)

		}
		sockets = append(sockets, socket)
		j++
	}

	return sockets
}

func (ts *TickersService) GetSocketsArrayLength() int {
	return ts.cache.PairsLength()
}

func (ts *TickersService) handleTransaction(tx *models.Transaction, index int) {
	if ts.writeToDB {
		ts.connections.InsertTransactionToDataBase(tx, index)
		ts.statsd.Client.Increment(".tickers.sqlinserts")

	} else {
		ts.cache.InsertTransaction(tx)
		ts.statsd.Client.Increment(".tickers.cacheinserts")

	}
}

func (ts *TickersService) handleOrderBookRow(tx *models.OrderBookRow, index int) {
	if ts.writeToDB {
		ts.connections.InsertOrderBookRowToDataBase(tx, index)
		ts.statsd.Client.Increment(".tickers.sqlinserts")

	} else {
		ts.cache.InsertOrderBookRow(tx)
		ts.statsd.Client.Increment(".tickers.cacheinserts")

	}
}

func (ts *TickersService) ListenAndHandle(txQueue chan *models.Transaction, obQueue chan *models.OrderBookRow, index int, wg *sync.WaitGroup, quit chan bool) {
	defer wg.Done()
	for {
		select {

		case <-quit:
			println("[ListenAndHandle] Exit")
			return

		case transaction := <-txQueue:
			ts.handleTransaction(transaction, index)

		case orderBookRow := <-obQueue:
			ts.handleOrderBookRow(orderBookRow, index)
		}

	}
}

func (ts *TickersService) GetTransactionChannel(index int) chan *models.Transaction {
	return ts.transactionChannels[index]
}

func (ts *TickersService) GetOrderBookChannel(index int) chan *models.OrderBookRow {
	return ts.orderBookChannels[index]
}

//TODO go routine grafana metric
