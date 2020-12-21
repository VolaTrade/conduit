package driver

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/service"
	"github.com/volatrade/conduit/internal/socket"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
)

var Module = wire.NewSet(
	New,
)

type (
	ConduitDriver struct {
		svc     service.Service
		kstats  *stats.Stats
		session *models.Session
		logger  *logger.Logger
	}
)

const (
	READS_PER_SECOND int = 5
)

type Driver interface {
	RunDataStreamListenerRoutines(wg *sync.WaitGroup, ch chan bool)
	RunSocketRecievingRoutines(wg *sync.WaitGroup)
	BuildDataChannels()
}

func New(svc service.Service, stats *stats.Stats, sess *models.Session, logger *logger.Logger) *ConduitDriver {
	return &ConduitDriver{svc: svc, kstats: stats, session: sess, logger: logger}
}

//BuildDataChannels ...
func (td *ConduitDriver) BuildDataChannels() {
	if err := td.svc.InsertPairsFromBinanceToCache(); err != nil {
		panic(err)
	}
	td.svc.BuildTransactionChannels(3)
	td.svc.BuildOrderBookChannels(3)

}

func (td *ConduitDriver) RunDataStreamListenerRoutines(wg *sync.WaitGroup, ch chan bool) {
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go td.svc.ListenAndHandleDataChannels(i, wg, ch)
	}
}

func (td *ConduitDriver) RunSocketRecievingRoutines(wg *sync.WaitGroup) {

	ctx, cancel := context.WithCancel(context.Background())

	go td.svc.ListenForDatabasePriveleges(wg)
	wg.Add(1)
	sockets := td.svc.SpawnSocketRoutines(3)
	go td.svc.ListenForExit(wg, cancel)
	go td.reportRunning(wg, ctx)

	for _, active_socket := range sockets {
		wg.Add(1)
		go td.consumeTransferTransactionMessage(active_socket, wg)
		go td.consumeTransferOrderBookMessage(active_socket, wg)
	}

}

func (td *ConduitDriver) consumeTransferTransactionMessage(socket *socket.BinanceSocket, wg *sync.WaitGroup) {
	defer wg.Done()
	td.logger.Infow("Consuming and transferring messsage")
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

func (cd *ConduitDriver) reportRunning(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	cd.kstats.Gauge(fmt.Sprintf("conduit.instances.%s", cd.session.ID), 1.0) //should this Gauge be 1?

	for {

		select {

		case <-ctx.Done():
			println("Reporting zero")
			cd.kstats.Gauge(fmt.Sprintf("conduit.instances.%s", cd.session.ID), 0.0)
			return

		}
	}
}
