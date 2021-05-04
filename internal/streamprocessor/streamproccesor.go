//go:generate mockgen -package=mocks -destination=../mocks/streamprocessor.go github.com/volatrade/conduit/internal/streamprocessor StreamProcessor
package streamprocessor

import (
	"context"
	"log"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/cache"
	"github.com/volatrade/conduit/internal/conveyor"
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/requests"
	"github.com/volatrade/conduit/internal/session"
	"github.com/volatrade/conduit/internal/storage"
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
		GenerateSocketListeningRoutines()
		GetProcessCollectionState() error
		ListenForDatabasePriveleges()
		ListenForExit(exit func())
		ListenAndHandleDataChannel(index int)
		RunSocketRoutines()
	}

	ConduitStreamProcessor struct {
		active            bool
		cache             cache.Cache
		conveyor          conveyor.Conveyor
		ctx               context.Context
		kstats            stats.Stats
		logger            *logger.Logger
		orderBookChannels []chan *models.OrderBookRow
		requests          requests.Requests
		slack             slack.Slack
		session           session.Session
	}
)

//New constructor
func New(ctx context.Context, conns storage.Store, ch cache.Cache, conveyor conveyor.Conveyor,
	cl requests.Requests, session session.Session, stats stats.Stats,
	slackz slack.Slack, logger *logger.Logger) (*ConduitStreamProcessor, func()) {

	sp := &ConduitStreamProcessor{
		ctx:      ctx,
		conveyor: conveyor,
		logger:   logger,
		cache:    ch,
		requests: cl,
		kstats:   stats,
		active:   false,
		slack:    slackz,
		session:  session,
	}

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
func (csp *ConduitStreamProcessor) handleOrderBookRow(ob *models.OrderBookRow, index int) {
	if csp.active {
		if err := csp.requests.PostOrderbookRowToCortex(ob); err != nil {
			csp.logger.Errorw("Error sending orderbook row to cortex", "error", err.Error())
		}
	}
	csp.cache.InsertOrderBookRow(ob)
	csp.kstats.Increment("cacheinserts.ob", 1.0)
}

//GetProcessCollectionState gathers the collection state for what pairs conduit should be collecting
func (csp *ConduitStreamProcessor) GetProcessCollectionState() error {

	tradingPairs, err := csp.requests.GetActiveOrderbookPairs(3)

	if err != nil {
		csp.logger.Errorw("Failed getting orderbook pairs from gatekeeper api, using default values")
		tradingPairs = []string{"btcusdt", "ethusdt", "xrpusdt", "ltcusdt"}
	}
	obChannelCount := int(len(tradingPairs) / 3)
	csp.BuildOrderBookChannels(obChannelCount)
	csp.logger.Infow("Fetching orderbook data", "pairs", tradingPairs, "channel count", obChannelCount)

	for _, pair := range tradingPairs {
		csp.cache.InsertEntry(pair)
	}

	return nil
}

//GetOrderBookChannel returns ob channel .... USED FOR TESTING ONLY
func (csp *ConduitStreamProcessor) GetOrderBookChannel(index int) chan *models.OrderBookRow {

	return csp.orderBookChannels[index]
}
