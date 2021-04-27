//go:generate mockgen -package=mocks -destination=../mocks/streamprocessor.go github.com/volatrade/conduit/internal/streamprocessor StreamProcessor
package streamprocessor

import (
	"context"
	"log"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/cache"
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
		BuildOrderBookChannels(count int)
		GenerateSocketListeningRoutines(ctx context.Context)
		InsertPairsFromBinanceToCache() error
		ListenForDatabasePriveleges(ctx context.Context)
		ListenForExit(exit func())
		ListenAndHandleDataChannel(ctx context.Context, index int)
		RunSocketRoutines(ctx context.Context)
	}

	ConduitStreamProcessor struct {
		logger            *logger.Logger
		cache             cache.Cache
		dbStreams         store.StorageConnections
		requests          requests.Requests
		slack             slack.Slack
		session           session.Session
		kstats            stats.Stats
		orderBookChannels []chan *models.OrderBookRow
		writeToDB         bool
	}
)

//New constructor
func New(conns store.StorageConnections, ch cache.Cache, cl requests.Requests, session session.Session, stats stats.Stats, slackz slack.Slack, logger *logger.Logger) (*ConduitStreamProcessor, func()) {

	sp := &ConduitStreamProcessor{
		logger:    logger,
		cache:     ch,
		dbStreams: conns,
		requests:  cl,
		kstats:    stats,
		writeToDB: false,
		slack:     slackz,
		session:   session,
	}

	sp.BuildOrderBookChannels(session.GetConnectionCount())

	shutdown := func() {
		log.Println("Shutting down stream proccessing layer")
		log.Println("Closing data channels")
		for i := 0; i < len(sp.orderBookChannels); i++ {
			close(sp.orderBookChannels[i])
		}
		log.Println("Stream processing layer completed shutdown")
	}
	return sp, shutdown
}

//handleOrderBookRow checks to see if orderbookrow is going to database or cache, then inserts accordingly
func (csp *ConduitStreamProcessor) handleOrderBookRow(tx *models.OrderBookRow, index int) {
	if csp.writeToDB {
		if err := csp.requests.PostOrderbookRowToCortex(tx); err != nil {
			csp.logger.Errorw("Error sending orderbook row to cortex", "error", err.Error())
		}
		if err := csp.dbStreams.InsertOrderBookRowToDataBase(tx, index); err != nil {
			csp.logger.Errorw("Error inserting orderbook row to postgres", "error", err.Error())
		}
		csp.kstats.Increment(".conduit.sqlinserts.ob", 1.0)

	} else {
		csp.cache.InsertOrderBookRow(tx)
		csp.kstats.Increment(".conduit.cacheinserts.ob", 1.0)

	}
}

//InsertPairsFromBinanceToCache reads all trading pairs from Binance and then proceeds to store them as keys in cache
func (csp *ConduitStreamProcessor) InsertPairsFromBinanceToCache() error {

	tradingPairs, err := csp.requests.GetActiveOrderbookPairs(3)

	if err != nil {
		csp.logger.Errorw("Failed getting orderbook pairs from gatekeeper api, using default values")
		tradingPairs = []string{"btcusdt", "ethusdt", "xrpusdt", "ltcusdt"}
	}

	csp.logger.Infow("Fetching orderbook data for: ", "pairs", tradingPairs)

	for _, pair := range tradingPairs {
		csp.cache.InsertEntry(pair)
	}

	return nil
}

//GetOrderBookChannel returns ob channel .... USED FOR TESTING ONLY
func (csp *ConduitStreamProcessor) GetOrderBookChannel(index int) chan *models.OrderBookRow {

	return csp.orderBookChannels[index]
}
