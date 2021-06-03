// conveyor package used to transfer mem cache on specified interval
package conveyor

import (
	"context"
	"time"

	"github.com/google/wire"

	"github.com/volatrade/conduit/internal/cache"
	"github.com/volatrade/conduit/internal/storage"
	logger "github.com/volatrade/currie-logs"
)

var Module = wire.NewSet(
	New,
)

type (
	Config struct {
		ShiftInterval int
	}
	Conveyor interface {
		Dispatch()
	}
	//ConduitConveyor ... kinda like a conveyor belt
	ConduitConveyor struct {
		cfg     *Config
		ctx     context.Context
		cache   cache.Cache
		logger  *logger.Logger
		storage storage.Store
	}
)

func New(cfg *Config, ctx context.Context, logger *logger.Logger,
	cache cache.Cache, storage storage.Store) *ConduitConveyor {
	return &ConduitConveyor{
		cfg:     cfg,
		ctx:     ctx,
		cache:   cache,
		storage: storage,
		logger:  logger,
	}
}

func (conv *ConduitConveyor) transitOrderBooksToStorage() {

	conv.logger.Infow("Starting cache to postgres transit operation")
	cachedObs := conv.cache.GetAllOrderBookRows()
	time.Sleep(300 * time.Millisecond)

	if len(cachedObs) != conv.cache.OrderBookRowsLength() { // checks to ensure this function isn't called amidst a cache update
		conv.logger.Infow("Recursive case hit within transit function")
		conv.transitOrderBooksToStorage()
	} else {
		conv.storage.TransferOrderBookCache(cachedObs)
		conv.cache.PurgeOrderBookRows()
	}
}

//Dispatch routine to wait for time interval criteria to be met before attempting cache to db migration
func (conv *ConduitConveyor) Dispatch() {
	ticker := time.NewTicker(time.Second * time.Duration(conv.cfg.ShiftInterval))
	defer ticker.Stop()
	conv.logger.Infow("Starting dispatch loop for conveyor")
	for {
		select {
		case <-ticker.C:
			if conv.cache.OrderBookRowsLength() != 0 {
				conv.transitOrderBooksToStorage()
			}
		case <-conv.ctx.Done():
			conv.logger.Infow("Received end signal from context in conveyor")
			return
		}
	}
}
