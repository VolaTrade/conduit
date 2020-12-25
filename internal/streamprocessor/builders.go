package streamprocessor

import (
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/socket"
)

//BuildTransactionChannels makes a slice of transaction struct channels
func (csp *ConduitStreamProcessor) BuildTransactionChannels(size int) {
	queues := make([]chan *models.Transaction, size)
	for i := 0; i < size; i++ {
		queue := make(chan *models.Transaction, 0)
		queues[i] = queue
	}
	csp.transactionChannels = queues
}

//BuildOrderBookChannels makes a slice of orderbook struct channels
func (csp *ConduitStreamProcessor) BuildOrderBookChannels(size int) {
	queues := make([]chan *models.OrderBookRow, size)

	for i := 0; i < size; i++ {
		queue := make(chan *models.OrderBookRow, 0)
		queues[i] = queue
	}

	csp.orderBookChannels = queues
}

//TODO add waitgroup to me ....
func (csp *ConduitStreamProcessor) RunSocketRoutines(psqlCount int) []*socket.ConduitSocketManager { // --> SpawnSocketManagers

	shepards := make([]*socket.ConduitSocketManager, 0)
	j := 0
	entries := csp.cache.GetEntries()
	for _, entry := range entries {
		if j >= psqlCount {
			j = 0
		}
		manager, err := socket.NewSocketManager(entry, csp.transactionChannels[j], csp.orderBookChannels[j], csp.kstats, csp.logger)

		if err != nil {
			csp.logger.Errorw(err.Error())

		}
		shepards = append(shepards, manager)
		j++
	}

	csp.logger.Infow("Spawned socket routines", "count", len(shepards))

	return shepards
}
