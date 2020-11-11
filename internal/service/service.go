package service

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

const (
	READS_PER_SECOND int = 5
)

type (
	Service interface {
		BuildPairUrls() error
		BuildTransactionChannels(count int)
		CheckForDatabasePriveleges()
		ChannelListenAndHandle(queue chan *models.Transaction, index int)
		ConsumeTransferMessage(socket *socket.BinanceSocket)
		SpawnSocketRoutines(psqlCount int) []*socket.BinanceSocket
		GetSocketsArrayLength() int
		GetChannel(index int) chan *models.Transaction
		ReportRunning()
	}

	TickersService struct {
		id                  string
		cache               cache.Cache
		connections         connections.Connections
		exch                client.Client
		statsd              *stats.StatsD
		transactionChannels []chan *models.Transaction
		writeToDB           bool
	}
)

func New(conns connections.Connections, ch cache.Cache, cl *client.ApiClient, stats *stats.StatsD) *TickersService {
	id := fmt.Sprintf("%d_%d", time.Now().Hour(), time.Now().Minute())
	return &TickersService{cache: ch, connections: conns, exch: cl, statsd: stats, writeToDB: false, id: id}
}

func (ts *TickersService) ReportRunning() {
	ts.statsd.Client.Gauge(fmt.Sprintf("tickers.instance.%s", ts.id), 1)
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

	for _, val := range tradingCryptosList {
		temp := val.(map[string]interface{}) //type casting
		id := strings.ToLower(temp["symbol"].(string))

		if strings.Contains(id, "btc") {
			ts.cache.InsertPairUrl(id)
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

func (ts *TickersService) SpawnSocketRoutines(psqlCount int) []*socket.BinanceSocket {

	sockets := make([]*socket.BinanceSocket, 0)
	j := 0
	for i := 0; i < ts.cache.PairUrlsLength()-1; i++ {
		if j >= psqlCount {
			j = 0
		}
		socketPairURL := ts.cache.GetPairUrl(i)
		temp_stats := stats.StatsD{}
		temp_stats.Client = ts.statsd.Client.Clone()
		pair := strings.Replace(socketPairURL, "wss://stream.binance.com:9443/ws/", "", 1)
		println("pair =---> ", pair)
		socket, err := socket.NewSocket(socketPairURL, pair, ts.transactionChannels[j], &temp_stats)
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
				ts.statsd.Client.Increment(".tickers.sqlinserts")

			} else {
				ts.cache.InsertTransaction(transaction)
				ts.statsd.Client.Increment(".tickers.cacheinserts")

			}
		}

	}
}

func (ts *TickersService) GetChannel(index int) chan *models.Transaction {
	return ts.transactionChannels[index]
}

//TODO go routine grafana metric

func (ts *TickersService) ConsumeTransferMessage(socket *socket.BinanceSocket) {

	var err error
	if err = socket.Connect(); err != nil {
		//TODO add handling policy

		panic(err)
	}
	secStored := int(time.Now().Second())
	hits := 0
	for {

		sec_now := time.Now().Second()
		if int(sec_now) > secStored || (secStored == 59 && sec_now != 59) {
			hits = 0
			secStored = sec_now
		}

		if hits >= READS_PER_SECOND {
			continue
		}
		message, err := socket.ReadMessage()

		if err != nil {
			//handle me
			log.Println(err.Error())
			ts.statsd.Client.Increment("tickers.errors.socket_read")
			continue
		}

		var transaction *models.Transaction

		if transaction, err = models.UnmarshalJSON(message); err != nil {
			ts.statsd.Client.Increment("tickers.errors.json_unmarshal")

			continue

		} else {
			log.Printf("%+v", transaction)
			socket.DataChannel <- transaction
		}
	}
}
