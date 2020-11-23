package dlq

import (
	"github.com/volatrade/tickers/internal/socket"
)

type (
	DeadQueue interface {
	}

	TickersDeadQueue struct {
		failedSockets []*socket.BinanceSocket
	}
)

func (tdq *TickersDeadQueue) InsertFailedSocket(socket *socket.BinanceSocket) {

	tdq.failedSockets = append(tdq.failedSockets, socket)
}
