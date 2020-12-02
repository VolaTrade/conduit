package driver

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/google/wire"
	"github.com/volatrade/tickers/internal/models"
	"github.com/volatrade/tickers/internal/service"
	"github.com/volatrade/tickers/internal/socket"
	"github.com/volatrade/tickers/internal/stats"
)

var Module = wire.NewSet(
	New,
)

type (
	TickersDriver struct {
		svc   service.Service
		statz *stats.StatsD
	}
)

const (
	READS_PER_SECOND int = 5
)

type Driver interface {
	RunListenerRoutines(wg *sync.WaitGroup, ch chan bool)
	Run(wg *sync.WaitGroup)
	InitService()
}

func New(svc service.Service, stats *stats.StatsD) *TickersDriver {
	return &TickersDriver{svc: svc, statz: stats}
}

//InitService initializes pairUrl list in cache and builds transactionChannels
func (td *TickersDriver) InitService() {
	if err := td.svc.BuildPairUrls(); err != nil {
		panic(err)
	}
	td.svc.BuildTransactionChannels(3)
	td.svc.BuildOrderBookChannels(3)

}

func (td *TickersDriver) RunListenerRoutines(wg *sync.WaitGroup, ch chan bool) {
	for i := 0; i < 3; i++ {
		wg.Add(1)
		txChannel := td.svc.GetTransactionChannel(i)
		obChannel := td.svc.GetOrderBookChannel(i)
		go td.svc.ListenAndHandle(txChannel, obChannel, i, wg, ch) //Tells channels to listen for transaction data from sockets
	}
}

func (td *TickersDriver) Run(wg *sync.WaitGroup) {
	go td.svc.CheckForDatabasePriveleges(wg)
	wg.Add(1)
	sockets := td.svc.SpawnSocketRoutines(3)
	go td.svc.ReportRunning(wg)
	for _, active_socket := range sockets {
		wg.Add(1)
		println("Spawning routine for -->", active_socket)
		go td.consumeTransferTransactionMessage(active_socket, wg)
		go td.consumeTransferOrderBookMessage(active_socket, wg)
		println("Spawned spawned")
	}

}

func (td *TickersDriver) consumeTransferTransactionMessage(socket *socket.BinanceSocket, wg *sync.WaitGroup) {
	defer wg.Done()
	println("Consuming and transferring messsage")
	var err error
	if err = socket.Connect(); err != nil {
		//TODO add handling policy
		println("error establishing socket connection")
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
			println("Continuing :p")
			continue
		}

		message, err := socket.ReadMessage("TX")

		if err != nil {
			//handle me
			log.Println(err.Error())
			td.statz.Client.Increment("tickers.errors.socket_read")
			continue
		}

		var transaction *models.Transaction

		if transaction, err = models.UnmarshalTransactionJSON(message); err != nil {
			println(err.Error())
			td.statz.Client.Increment("tickers.errors.json_unmarshal")

		} else {
			log.Printf("%+v", transaction)
			socket.TransactionChannel <- transaction
		}

		// TODO: Add order book insertions
		// TODO: Add support for passing pair since we dont get it back from socket
	}
}

func getCurrentMilisecond() (int, error) {
	ts := int(time.Now().UnixNano()) / (int(time.Millisecond) / int(time.Nanosecond))

	strTime := strconv.Itoa(ts)

	ms, err := strconv.Atoi(strTime[10:11])

	return ms, err
}
func (td *TickersDriver) consumeTransferOrderBookMessage(socket *socket.BinanceSocket, wg *sync.WaitGroup) {
	defer wg.Done()
	println("Consuming and transferring messsage")
	var err error
	if err = socket.Connect(); err != nil {
		//TODO add handling policy
		println("error establishing socket connection")
		panic(err)
	}
	reset := true
	for {
		ms, _ := getCurrentMilisecond()

		if ms%2 != 0 && ms != 0 {
			reset = true
			continue
		}

		if !reset {
			continue
		}

		message, err := socket.ReadMessage("OB")

		reset = false
		println("RAW ORDER BOOK MESSAGE ->", string(message))
		println("READ @ ->", ms)
		if err != nil {
			//handle me
			log.Println(err.Error())
			td.statz.Client.Increment("tickers.errors.socket_read")
			continue
		}

		var orderBookRow *models.OrderBookRow

		if orderBookRow, err = models.UnmarshalOrderBookJSON(message, socket.Pair); err != nil {
			println(err.Error())
			td.statz.Client.Increment("tickers.errors.json_unmarshal")

		} else {
			log.Printf("%+v", orderBookRow)
			socket.OrderBookChannel <- orderBookRow
		}
	}
}
