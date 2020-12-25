package socket

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	logger "github.com/volatrade/currie-logs"
)

/*
TODO
1. Add context to me
2. unit test me
3. add select for context channel
*/

const (
	TIMEOUT = time.Second
)

type ConduitSocket struct {
	parentChannel chan bool
	mux           *sync.Mutex
	logger        *logger.Logger
	conn          *websocket.Conn
	url           string
	healthy       bool
}

func NewSocket(url string, logger *logger.Logger, channel chan bool) (*ConduitSocket, error) {

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		return nil, err
	}

	socket := &ConduitSocket{conn: conn, url: url, logger: logger, parentChannel: channel, healthy: true, mux: &sync.Mutex{}}

	go socket.runKeepAlive()

	return socket, nil

}

func (cs *ConduitSocket) readMessage() ([]byte, error) {

	cs.mux.Lock()
	defer cs.mux.Unlock()
	cs.conn.SetReadDeadline(time.Now().Add(TIMEOUT))

	_, message, err := cs.conn.ReadMessage()

	return message, err

}

func (cs *ConduitSocket) runKeepAlive() {

	cs.logger.Infow("keep alive", "url", cs.url)
	ticker := time.NewTicker(time.Second * 15)
	for {

		if _, err := cs.readMessage(); err != nil || !cs.healthy {
			cs.healthy = false
			cs.logger.Errorw(err.Error(), "url", cs.url, "n")
			go cs.reconnect(3)
			return
		}
		select {

		case <-cs.parentChannel:
			cs.healthy = false
			go cs.reconnect(3)
			return

		case <-ticker.C: //Ticker interval triggered
			continue

		}
	}
}

func (cs *ConduitSocket) reconnect(times int) {
	cs.conn.Close()

	if times == 0 {
		cs.logger.Errorw("reconnection failed ... broken socket ... ")
		cs.parentChannel <- false
	}

	cs.logger.Infow("attempting to reconnect to failed socket", "attempt", times)

	cs.mux.Lock()
	conn, _, err := websocket.DefaultDialer.Dial(cs.url, nil)

	if err != nil {
		cs.logger.Infow("reconnection attempt failed", "attempt", times)
		cs.mux.Unlock()
		cs.reconnect(times - 1)

	}
	cs.conn = conn
	cs.mux.Unlock()
	cs.logger.Infow("successfully restarted failed socket")
	cs.healthy = true
	cs.runKeepAlive()

}
