package streamprocessor

import (
	"os"
	"time"
)

//ListenAndHandleDataChannel waits for transaction or orderbook to come in from their respective channel before invoking handler function
func (csp *ConduitStreamProcessor) ListenAndHandleDataChannel(index int) {

	obChannel := csp.orderBookChannels[index]
	for {
		select {

		case orderBookRow := <-obChannel:
			csp.handleOrderBookRow(orderBookRow, index)

		case <-csp.ctx.Done():
			csp.logger.Infow("received finish signal")
			return
		}

	}
}

//ListenForDatabasePriveleges checks every fifteen seconds for prevalance of touch file, signals database writing when found
func (csp *ConduitStreamProcessor) ListenForDatabasePriveleges() {

	ticker := time.NewTicker(time.Duration(time.Second * 15))
	defer ticker.Stop()
	attempts := 3

	for {
		if attempts == 0 {
			csp.logger.Errorw("Failed to perform database operations")
			panic("FAILURE")
		}

		if _, startDetected := os.Stat("start"); startDetected == nil {

			csp.logger.Infow("Dispatching conveyor")
			csp.active = true
			go csp.conveyor.Dispatch()
			return
		}

		select {
		case <-ticker.C:
			continue

		case <-csp.ctx.Done():
			csp.logger.Infow("received finish signal")
			return

		}
	}
}

//ListenForExit checks every fifteen seconds for prevalance of finish file
func (csp *ConduitStreamProcessor) ListenForExit(exit func()) {

	ticker := time.NewTicker(time.Duration(time.Second * 15))
	defer ticker.Stop()
	for {
		if _, err := os.Stat("finish"); err == nil {
			csp.logger.Infow("Finish signal recieved")
			exit()
			return
		}

		for range ticker.C {
			break
		}
	}
}

func (csp *ConduitStreamProcessor) GenerateSocketListeningRoutines() {
	for i := 0; i < len(csp.orderBookChannels); i++ {
		go csp.ListenAndHandleDataChannel(i)
	}
}
