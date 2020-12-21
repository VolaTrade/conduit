//go:generate mockgen -package=mocks -destination=../mocks/service.go github.com/volatrade/conduit/internal/service Service
package service

import (
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
	Service interface {
		BuildTransactionChannels(count int)
		BuildOrderBookChannels(count int)
		InsertPairsFromBinanceToCache() error
		ListenForDatabasePriveleges(wg *sync.WaitGroup)
		ListenForExit(wg *sync.WaitGroup, exit func())
		ListenAndHandleDataChannels(index int, wg *sync.WaitGroup, ch chan bool)
		SpawnSocketRoutines(psqlCount int) []*socket.BinanceSocket
	}

	ConduitService struct {
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

func New(conns store.StorageConnections, ch cache.Cache, cl requests.Requests, stats *stats.Stats, slackz slack.Slack, logger *logger.Logger) *ConduitService {

	return &ConduitService{
		logger:    logger,
		cache:     ch,
		dbStreams: conns,
		requests:  cl,
		kstats:    stats,
		writeToDB: false,
		slack:     slackz,
	}
}

func (ts *ConduitService) handleTransaction(tx *models.Transaction, index int) {
	if ts.writeToDB {
		ts.dbStreams.InsertTransactionToDataBase(tx, index)
		ts.kstats.Increment(".conduit.sqlinserts.tx", 1.0)

	} else {
		ts.cache.InsertTransaction(tx)
		ts.kstats.Increment(".conduit.cacheinserts.tx", 1.0)

	}
}

func (ts *ConduitService) handleOrderBookRow(tx *models.OrderBookRow, index int) {
	if ts.writeToDB {
		ts.dbStreams.InsertOrderBookRowToDataBase(tx, index)
		ts.kstats.Increment(".conduit.sqlinserts.ob", 1.0)

	} else {
		ts.cache.InsertOrderBookRow(tx)
		ts.kstats.Increment(".conduit.cacheinserts.ob", 1.0)

	}
}

//InsertPairsFromBinanceToCache reads all trading pairs from Binance and then proceeds to store them as keys in cache
func (ts *ConduitService) InsertPairsFromBinanceToCache() error {

	tradingPairs, err := ts.requests.GetActiveBinanceExchangePairs()

	if err != nil {
		ts.logger.Errorw(err.Error())
		return err
	}

	for _, pair := range tradingPairs {

		if pair == "btcusdt" || pair == "ethusdt" || pair == "xrpusdt" {
			ts.cache.InsertEntry(pair)
		}
	}

	return nil
}
