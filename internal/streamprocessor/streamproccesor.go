//go:generate mockgen -package=mocks -destination=../mocks/streamprocessor.go github.com/volatrade/conduit/internal/streamprocessor StreamProcessor
package streamprocessor

import (
	"context"
	"log"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/cache"
	"github.com/volatrade/conduit/internal/cortex"
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/requests"
	"github.com/volatrade/conduit/internal/session"
	"github.com/volatrade/conduit/internal/store"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
	"github.com/volatrade/utilities/slack"
)

var (
	Module = wire.NewSet(
		New,
	)
)

type (
	StreamProcessor interface {
		BuildTransactionChannels(count int)
		BuildOrderBookChannels(count int)
		GenerateSocketListeningRoutines(ctx context.Context)
		InsertPairsFromBinanceToCache() error
		ListenForDatabasePriveleges(ctx context.Context)
		ListenForExit(exit func())
		ListenAndHandleDataChannel(ctx context.Context, index int)
		RunSocketRoutines(ctx context.Context)
	}

	ConduitStreamProcessor struct {
		cortex              cortex.Cortex
		logger              *logger.Logger
		cache               cache.Cache
		dbStreams           store.StorageConnections
		requests            requests.Requests
		slack               slack.Slack
		session             session.Session
		kstats              *stats.Stats
		transactionChannels []chan *models.Transaction
		orderBookChannels   []chan *models.OrderBookRow
		writeToDB           bool
	}
)

//New constructor
func New(conns store.StorageConnections, ch cache.Cache, cl requests.Requests, session session.Session, stats *stats.Stats, slackz slack.Slack, logger *logger.Logger, cortex cortex.Cortex) (*ConduitStreamProcessor, func()) {

	sp := &ConduitStreamProcessor{
		logger:    logger,
		cache:     ch,
		dbStreams: conns,
		requests:  cl,
		kstats:    stats,
		writeToDB: false,
		slack:     slackz,
		session:   session,
		cortex:    cortex,
	}

	sp.BuildTransactionChannels(session.GetConnectionCount())
	sp.BuildOrderBookChannels(session.GetConnectionCount())

	shutdown := func() {
		log.Println("Shutting down stream proccessing layer")
		log.Println("Closing data channels")
		for i := 0; i < len(sp.transactionChannels); i++ {
			close(sp.transactionChannels[i])
			close(sp.orderBookChannels[i])
		}
		log.Println("Stream processing layer completed shutdown")
	}
	return sp, shutdown
}

//handleTransaction checks to see if tx is going to database or cache, then inserts accordingly
func (csp *ConduitStreamProcessor) handleTransaction(tx *models.Transaction, index int) {
	if csp.writeToDB {
		csp.dbStreams.InsertTransactionToDataBase(tx, index)
		csp.kstats.Increment(".conduit.sqlinserts.tx", 1.0)

	} else {
		csp.cache.InsertTransaction(tx)
		csp.kstats.Increment(".conduit.cacheinserts.tx", 1.0)

	}
}

//handleOrderBookRow checks to see if orderbookrow is going to database or cache, then inserts accordingly
func (csp *ConduitStreamProcessor) handleOrderBookRow(ob *models.OrderBookRow, index int) {
	//send data to cortex

	csp.cortex.SendOrderBookRow(ob)

	if csp.writeToDB {
		csp.dbStreams.InsertOrderBookRowToDataBase(ob, index)
		csp.kstats.Increment(".conduit.sqlinserts.ob", 1.0)

	} else {
		csp.cache.InsertOrderBookRow(ob)
		csp.kstats.Increment(".conduit.cacheinserts.ob", 1.0)
	}
}

//InsertPairsFromBinanceToCache reads all trading pairs from Binance and then proceeds to store them as keys in cache
func (csp *ConduitStreamProcessor) InsertPairsFromBinanceToCache() error {

	tradingPairs, err := csp.requests.GetActiveBinanceExchangePairs()

	if err != nil {
		csp.logger.Errorw(err.Error())
		return err
	}

	for _, pair := range tradingPairs {

		if pair == "btcusdt" || pair == "ethusdt" || pair == "xrpusdt" || pair == "ltcusdt" {
			csp.cache.InsertEntry(pair)
		}
	}

	return nil
}

//GetOrderBookChannel returns ob channel .... USED FOR TESTING ONLY
func (csp *ConduitStreamProcessor) GetOrderBookChannel(index int) chan *models.OrderBookRow {

	return csp.orderBookChannels[index]
}

//GetOrderBookChannel returns tx channel .... USED FOR TESTING ONLY
func (csp *ConduitStreamProcessor) GetTransactionChannel(index int) chan *models.Transaction {
	return csp.transactionChannels[index]
}
