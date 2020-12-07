package dlq

import (
	"github.com/volatrade/conduit/internal/socket"
)

type (
	DeadQueue interface {
	}

	conduitDeadQueue struct {
		failedSockets []*socket.BinanceSocket
	}
)

func (tdq *conduitDeadQueue) InsertFailedSocket(socket *socket.BinanceSocket) {

	tdq.failedSockets = append(tdq.failedSockets, socket)
}
