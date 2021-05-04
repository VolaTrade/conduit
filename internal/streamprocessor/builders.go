package streamprocessor

import (
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/socket"
)

//BuildOrderBookChannels makes a slice of orderbook struct channels
func (csp *ConduitStreamProcessor) BuildOrderBookChannels(size int) {
	queues := make([]chan *models.OrderBookRow, size)

	for i := 0; i < size; i++ {
		queue := make(chan *models.OrderBookRow)
		queues[i] = queue
	}

	csp.orderBookChannels = queues
}

//TODO add waitgroup to me ....
func (csp *ConduitStreamProcessor) RunSocketRoutines() { // --> SpawnSocketManagers

	shepards := make([]*socket.ConduitSocketManager, 0)
	j := 0
	entries := csp.cache.GetEntries()
	connCount := len(csp.orderBookChannels)
	for _, entry := range entries {
		if j >= connCount {
			j = 0
		}
		manager := socket.NewSocketManager(entry, csp.orderBookChannels[j], csp.kstats, csp.logger)

		shepards = append(shepards, manager)
		j++
	}

	csp.logger.Infow("Built socket routines", "count", len(shepards))

	for _, manager := range shepards {
		go manager.Run(csp.ctx)
	}
}
