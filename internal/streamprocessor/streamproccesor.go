//go:generate mockgen -package=mocks -destination=../mocks/streamprocessor.go github.com/volatrade/conduit/internal/streamprocessor StreamProcessor
package streamprocessor

import (
	"context"
	"log"
	"time"

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
		cortexClient      cortex.Cortex
		orderBookChannels []chan *models.OrderBookRow
		writeToDB         bool
	}
)

//New constructor
func New(conns store.StorageConnections, ch cache.Cache, cl requests.Requests, session session.Session,
	stats stats.Stats, slackz slack.Slack, logger *logger.Logger, cortexClient cortex.Cortex) (*ConduitStreamProcessor, func()) {

	sp := &ConduitStreamProcessor{
		logger:       logger,
		cache:        ch,
		dbStreams:    conns,
		requests:     cl,
		kstats:       stats,
		writeToDB:    false,
		slack:        slackz,
		session:      session,
		cortexClient: cortexClient,
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

func (csp *ConduitStreamProcessor) ProcessObRowsToCortex(ob *models.OrderBookRow) error {

	if err := csp.cache.InsertOrderBookRowToRedis(ob); err != nil {
		return err
	}

	obRows, err := csp.cache.GetOrderBookRowsFromRedis(ob.Pair)

	if err != nil {
		csp.kstats.Increment("conduit.sent_obrow.cortex.error", 1.0)
		return err
	}

	start := time.Now()
	defer csp.kstats.TimingDuration("conduit.sent_obrow.cortex.time_duration", time.Since(start))

	if err := csp.cortexClient.SendOrderBookRows(obRows); err != nil {
		csp.logger.Errorw(err.Error(), "error sending message")
		csp.kstats.Increment("conduit.sent_obrow.cortex.error", 1.0)
		return err
	}

	return nil
}

//handleOrderBookRow checks to see if orderbookrow is going to database or cache, then inserts accordingly
func (csp *ConduitStreamProcessor) handleOrderBookRow(ob *models.OrderBookRow, index int) {
	start := time.Now()
	defer csp.kstats.TimingDuration("conduit.orderbook_row_handling.time_duration", time.Since(start))
	defer csp.logger.Infow("Orderbook handling logic TTL", "time", time.Since(start).String())

	if csp.cache.RowValidForCortex(ob.Pair) {

		if err := csp.ProcessObRowsToCortex(ob); err != nil {
			csp.logger.Errorw(err.Error())
		}

	}

	if csp.writeToDB {
		csp.dbStreams.InsertOrderBookRowToDataBase(ob, index)
		csp.kstats.Increment("conduit.sqlinserts.ob", 1.0)

	} else {
		csp.cache.InsertOrderBookRow(ob)
		csp.kstats.Increment("conduit.cacheinserts.ob", 1.0)

	}
}

//InsertPairsFromBinanceToCache reads all trading pairs from Binance and then proceeds to store them as keys in cache
func (csp *ConduitStreamProcessor) InsertPairsFromBinanceToCache() error {

	tradingPairs := []string{"btcusdt", "ethusdt", "xrpusdt", "ltcusdt"}

	for _, pair := range tradingPairs {
		csp.cache.InsertEntry(pair)
	}

	return nil
}

//GetOrderBookChannel returns ob channel .... USED FOR TESTING ONLY
func (csp *ConduitStreamProcessor) GetOrderBookChannel(index int) chan *models.OrderBookRow {

	return csp.orderBookChannels[index]
}
