package service

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/wire"
	"github.com/volatrade/tickers/internal/cache"
	"github.com/volatrade/tickers/internal/client"
	"github.com/volatrade/tickers/internal/connections"
	"github.com/volatrade/tickers/internal/models"
	"github.com/volatrade/tickers/internal/socket"
	"github.com/volatrade/tickers/internal/stats"
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
		CheckForDatabasePriveleges()
		ChannelListenAndHandle(queue chan *models.Transaction, index int)
		SpawnSocketRoutines(psqlCount int) []*socket.BinanceSocket
		GetSocketsArrayLength() int
		GetChannel(index int) chan *models.Transaction
	}

	TickersService struct {
		cache               cache.Cache
		connections         connections.Connections
		exch                client.Client
		statsd              *stats.StatsD
		transactionChannels []chan *models.Transaction
		writeToDB           bool
	}
)

func New(conns connections.Connections, ch cache.Cache, cl *client.ApiClient, stats *stats.StatsD) *TickersService {
	return &TickersService{cache: ch, connections: conns, exch: cl, statsd: stats, writeToDB: false}
}

//TODO there's a better way to structure this
func (ts *TickersService) CheckForDatabasePriveleges() {

	for {

		if _, err := os.Stat("start"); err == nil {
			log.Println("making connections to DB NOW")
			ts.connections.MakeConnections()
			ts.writeToDB = true
			log.Println("Purging cache")
			cached_transactions := ts.cache.GetAllTransactions()
			err := ts.connections.TransferCache(cached_transactions)
			if err != nil {
				panic(err)
			}
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

	for i, val := range tradingCryptosList {
		temp := val.(map[string]interface{}) //type casting
		id := strings.ToLower(temp["symbol"].(string))

		if strings.Contains(id, "btc") {
			ts.cache.InsertPairUrl(id)
		}
		if i >= 50 {
			break
		}
	}

	return nil
}

func (ts *TickersService) BuildTransactionChannels(count int) {
	queues := make([]chan *models.Transaction, count)
	for i := 0; i < count; i++ {
		queue := make(chan *models.Transaction, 1)
		queues[i] = queue
	}
	ts.transactionChannels = queues
}

func (ts *TickersService) SpawnSocketRoutines(psqlCount int) []*socket.BinanceSocket {

	sockets := make([]*socket.BinanceSocket, 0)
	j := 0
	for i := 0; i < ts.cache.PairUrlsLength()-1; i++ {
		if j >= psqlCount {
			j = 0
		}
		println("getting pair url for -->", ts.cache.GetPairUrl(i))
		socketPairURL := ts.cache.GetPairUrl(i)
		socket, err := socket.NewSocket(socketPairURL, ts.transactionChannels[j])
		println("Built socket for -->", socketPairURL)
		if err != nil {
			fmt.Println(err)

		}
		sockets = append(sockets, socket)
		j++
	}

	return sockets

}
func (ts *TickersService) GetSocketsArrayLength() int {
	return ts.cache.PairUrlsLength()
}
func (ts *TickersService) ChannelListenAndHandle(queue chan *models.Transaction, index int) {

	for {
		for transaction := range queue {

			if ts.writeToDB {
				ts.connections.InsertTransactionToDataBase(transaction, index)
				ts.statsd.Client.Increment(".transacts.sqlinserts")

			} else {
				ts.cache.InsertTransaction(transaction)
				ts.statsd.Client.Increment(".transacts.cacheinserts")

			}
		}

	}
}

func (ts *TickersService) GetChannel(index int) chan *models.Transaction {
	return ts.transactionChannels[index]
}

//TODO go routine grafana metric
