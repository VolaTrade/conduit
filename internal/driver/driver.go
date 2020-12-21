package driver

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/service"
	"github.com/volatrade/conduit/internal/socket"
	stats "github.com/volatrade/k-stats"
)

var Module = wire.NewSet(
	New,
)

type (
	ConduitDriver struct {
		svc    service.Service
		kstats *stats.Stats
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

func New(svc service.Service, stats *stats.Stats) *ConduitDriver {
	return &ConduitDriver{svc: svc, kstats: stats}
}

//InitService initializes pairUrl list in cache and builds transactionChannels
func (td *ConduitDriver) InitService() {
	if err := td.svc.BuildPairUrls(); err != nil {
		panic(err)
	}
	td.svc.BuildTransactionChannels(3)
	td.svc.BuildOrderBookChannels(3)

}

func (td *ConduitDriver) RunListenerRoutines(wg *sync.WaitGroup, ch chan bool) {
	for i := 0; i < 3; i++ {
		wg.Add(1)
		txChannel := td.svc.GetTransactionChannel(i)
		obChannel := td.svc.GetOrderBookChannel(i)
		go td.svc.ListenAndHandle(txChannel, obChannel, i, wg, ch) //Tells channels to listen for transaction data from sockets
	}
}

func (td *ConduitDriver) Run(wg *sync.WaitGroup) {

	ctx, cancel := context.WithCancel(context.Background())

	go td.svc.CheckForDatabasePriveleges(wg)
	wg.Add(1)
	sockets := td.svc.SpawnSocketRoutines(3)
	go td.svc.ReportRunning(wg, ctx)
	go td.svc.CheckForExit(wg, cancel)
	for _, active_socket := range sockets {
		wg.Add(1)
		log.Println("Spawning routine for -->", active_socket)
		go td.consumeTransferTransactionMessage(active_socket, wg)
		go td.consumeTransferOrderBookMessage(active_socket, wg)
		log.Println("Spawned spawned")
	}

}

func (td *ConduitDriver) consumeTransferTransactionMessage(socket *socket.BinanceSocket, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Consuming and transferring messsage")
	var err error
	if err = socket.Connect(); err != nil {
		//TODO add handling policy
		log.Println("error establishing socket connection")
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

		message, err := socket.ReadMessage("TX")

		if err != nil {
			//handle me
			log.Println(err.Error())
			td.kstats.Increment("conduit.errors.socket_read.tx", 1.0)
			continue
		}

		var transaction *models.Transaction

		if transaction, err = models.UnmarshalTransactionJSON(message); err != nil {
			log.Println(err.Error())
			td.kstats.Increment("conduit.errors.json_unmarshal", 1.0)

		} else {
			socket.TransactionChannel <- transaction
		}

		// TODO: Add order book insertions
		// TODO: Add support for passing pair since we dont get it back from socket
	}
}

func runKeepAlive(socket *socket.BinanceSocket, kstats *stats.Stats, mux *sync.Mutex) {
	prev_sec := time.Now().Second()

	for {
		curr_sec := time.Now().Second()

		if prev_sec == curr_sec {
			continue
		}
		mux.Lock()
		_, err := socket.ReadMessage("OB")
		mux.Unlock()
		prev_sec = curr_sec

		if err != nil {
			kstats.Increment("conduit.errors.socket_read.ob", 1.0)
		}

	}
}

func (td *ConduitDriver) consumeTransferOrderBookMessage(socket *socket.BinanceSocket, wg *sync.WaitGroup) {
	socketMux := &sync.Mutex{}
	defer wg.Done()
	log.Println("Consuming and transferring messsage")
	var err error
	if err = socket.Connect(); err != nil {
		//TODO add handling policy
		log.Println("error establishing socket connection", 1.0)
		panic(err)
	}

	go func() {
		runKeepAlive(socket, td.kstats, socketMux)
	}()

	prev_min := time.Now().Minute() - 1
	for {

		curr_min := time.Now().Minute()

		if prev_min == curr_min {
			continue
		}
		socketMux.Lock()
		message, err := socket.ReadMessage("OB")
		socketMux.Unlock()
		prev_min = curr_min
		if err != nil {
			//handle me
			log.Println(err.Error(), socket.Pair)
			td.kstats.Increment("conduit.errors.socket_read.ob", 1.0)
			continue
		}

		var orderBookRow *models.OrderBookRow

		if orderBookRow, err = models.UnmarshalOrderBookJSON(message, socket.Pair); err != nil {
			log.Println(err.Error(), socket.Pair)
			td.kstats.Increment("conduit.errors.json_unmarshal", 1.0)

		} else {
			socket.OrderBookChannel <- orderBookRow
		}
	}
}
