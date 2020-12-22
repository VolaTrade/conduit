package service

import (
	"github.com/volatrade/conduit/internal/models"
	"github.com/volatrade/conduit/internal/socket"
)

//BuildTransactionChannels makes a slice of transaction struct channels
func (ts *ConduitService) BuildTransactionChannels(size int) {
	queues := make([]chan *models.Transaction, size)
	for i := 0; i < size; i++ {
		queue := make(chan *models.Transaction, 0)
		queues[i] = queue
	}
	ts.transactionChannels = queues
}

//BuildOrderBookChannels makes a slice of orderbook struct channels
func (ts *ConduitService) BuildOrderBookChannels(size int) {
	queues := make([]chan *models.OrderBookRow, size)

	for i := 0; i < size; i++ {
		queue := make(chan *models.OrderBookRow, 0)
		queues[i] = queue
	}

	ts.orderBookChannels = queues
}

func (ts *ConduitService) SpawnSocketRoutines(psqlCount int) []*socket.BinanceSocket { // --> SpawnSocketManagers

	sockets := make([]*socket.BinanceSocket, 0)
	j := 0
	entries := ts.cache.GetEntries()
	for _, entry := range entries {
		if j >= psqlCount {
			j = 0
		}
		socket, err := socket.NewSocket(entry.TxUrl, entry.ObUrl, entry.Pair, ts.transactionChannels[j], ts.orderBookChannels[j], ts.kstats, ts.logger) //socket --> socket manager

		if err != nil {
			ts.logger.Errorw(err.Error())

		}
		sockets = append(sockets, socket)
		j++
	}

	ts.logger.Infow("Spawned socket routines", "count", len(sockets))

	return sockets
}