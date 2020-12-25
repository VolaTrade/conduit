//go:generate mockgen -package=mocks -destination=../mocks/streamprocessor.go github.com/volatrade/conduit/internal/streamprocessor StreamProcessor
package streamprocessor

import (
	"context"
	"log"
	"sync"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/cache"
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/requests"
	"github.com/volatrade/conduit/internal/socket"
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
		InsertPairsFromBinanceToCache() error
		ListenForDatabasePriveleges(ctx context.Context, wg *sync.WaitGroup)
		ListenForExit(ctx context.Context, wg *sync.WaitGroup, exit func())
		ListenAndHandleDataChannels(ctx context.Context, index int, wg *sync.WaitGroup)
		RunSocketRoutines(psqlCount int) []*socket.ConduitSocketManager
	}

	ConduitStreamProcessor struct {
		logger              *logger.Logger
		cache               cache.Cache
		dbStreams           store.StorageConnections
		requests            requests.Requests
		slack               slack.Slack
		kstats              *stats.Stats
		transactionChannels []chan *models.Transaction
		orderBookChannels   []chan *models.OrderBookRow
		writeToDB           bool
	}
)

//New constructor
func New(conns store.StorageConnections, ch cache.Cache, cl requests.Requests, stats *stats.Stats, slackz slack.Slack, logger *logger.Logger) (*ConduitStreamProcessor, func()) {

	sp := &ConduitStreamProcessor{
		logger:    logger,
		cache:     ch,
		dbStreams: conns,
		requests:  cl,
		kstats:    stats,
		writeToDB: false,
		slack:     slackz,
	}

	sp.BuildTransactionChannels(3)
	sp.BuildOrderBookChannels(3)

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
func (csp *ConduitStreamProcessor) handleOrderBookRow(tx *models.OrderBookRow, index int) {
	if csp.writeToDB {
		csp.dbStreams.InsertOrderBookRowToDataBase(tx, index)
		csp.kstats.Increment(".conduit.sqlinserts.ob", 1.0)

	} else {
		csp.cache.InsertOrderBookRow(tx)
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

		if pair == "btcusdt" || pair == "ethusdt" || pair == "xrpusdt" {
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
