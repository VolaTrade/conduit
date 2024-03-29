package socket

import (
	"sync"
	"time"

	"context"

	"github.com/gorilla/websocket"
	logger "github.com/volatrade/currie-logs"
)

/*
TODO
unit test me
*/

const (
	TIMEOUT = time.Second * 2
)

type ConduitSocket struct {
	parentChannel chan bool
	mux           *sync.Mutex
	logger        *logger.Logger
	conn          *websocket.Conn
	url           string
	healthy       bool
	ctx           context.Context
}

func NewSocket(ctx context.Context, url string, logger *logger.Logger, channel chan bool) (*ConduitSocket, error) {

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		return nil, err
	}
	socket := &ConduitSocket{conn: conn, url: url, logger: logger, parentChannel: channel, healthy: true, mux: &sync.Mutex{}, ctx: ctx}

	go socket.runKeepAlive()

	return socket, nil

}

func (cs *ConduitSocket) readMessage() ([]byte, error) {
	// start := time.Now()
	cs.mux.Lock()
	//TODO make this grafana panel
	// defer cs.stas.Infow("read message complete", "time", time.Since(start), "url", cs.url)
	defer cs.mux.Unlock()

	if err := cs.conn.SetReadDeadline(time.Now().Add(TIMEOUT)); err != nil {
		return []byte{}, err
	}

	_, message, err := cs.conn.ReadMessage()

	return message, err

}

func (cs *ConduitSocket) runKeepAlive() {

	cs.logger.Infow("keep alive", "url", cs.url)
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
	for {

		// start := time.Now()
		if err := cs.conn.WriteMessage(websocket.PongMessage, []byte("")); err != nil || !cs.healthy {
			cs.healthy = false
			cs.logger.Errorw(err.Error(), "url", cs.url)
			go cs.reconnect(3)
			return
		}

		// TODO make this a grafana metric
		// cs.logger.Infow("write message complete", "time", time.Since(start), "url", cs.url)
		select {

		case <-cs.parentChannel:
			cs.healthy = false
			go cs.reconnect(3)
			return

		case <-ticker.C: //Ticker interval triggered
			continue

		case <-cs.ctx.Done():
			cs.logger.Infow("received finish sig from context", "url", cs.url)
			return
		}
	}
}

func (cs *ConduitSocket) reconnect(times int) {
	cs.conn.Close()

	if times == 0 {
		cs.logger.Errorw("reconnection failed ... broken socket ... ")
		cs.parentChannel <- false
	}

	cs.logger.Infow("attempting to reconnect to failed socket", "attempt", times, "url", cs.url)

	cs.mux.Lock()
	conn, _, err := websocket.DefaultDialer.Dial(cs.url, nil)

	if err != nil {
		cs.logger.Infow("reconnection attempt failed", "attempt", times)
		cs.mux.Unlock()
		cs.reconnect(times - 1)

	}
	cs.conn = conn
	cs.mux.Unlock()
	cs.logger.Infow("successfully restarted failed socket", "url", cs.url)
	cs.healthy = true
	cs.runKeepAlive()

}
