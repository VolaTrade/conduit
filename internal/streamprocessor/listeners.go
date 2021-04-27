package streamprocessor

import (
	"context"
	"os"
	"time"
)

//ListenAndHandleDataChannel waits for transaction or orderbook to come in from their respective channel before invoking handler function
func (csp *ConduitStreamProcessor) ListenAndHandleDataChannel(ctx context.Context, index int) {

	obChannel := csp.orderBookChannels[index]
	for {
		select {

		case orderBookRow := <-obChannel:
			csp.handleOrderBookRow(orderBookRow, index)

		case <-ctx.Done():
			csp.logger.Infow("received finish signal")
			return
		}

	}
}

//ListenForDatabasePriveleges checks every fifteen seconds for prevalance of touch file, signals database writing when found
func (csp *ConduitStreamProcessor) ListenForDatabasePriveleges(ctx context.Context) {

	ticker := time.NewTicker(time.Duration(time.Second * 15))
	defer ticker.Stop()
	attempts := 3

	for {
		if attempts == 0 {
			csp.logger.Errorw("Failed to perform database operations")
			panic("FAILURE")
		}

		if _, writeToCache := os.Stat("start"); writeToCache == nil {

			csp.logger.Infow("establishing database connections, moving cache to databse, and purging cache")

			if err := csp.dbStreams.MakeConnections(); err != nil {
				csp.logger.Errorw(err.Error(), "attempt", attempts)
				attempts -= 1
				continue
			}

			csp.writeToDB = true
			if err := csp.dbStreams.TransferOrderBookCache(csp.cache.GetAllOrderBookRows()); err != nil {
				csp.logger.Errorw(err.Error(), "attempt", attempts)
				attempts -= 1
				continue
			}

			csp.cache.PurgeOrderBookRows()
			return
		}

		select {
		case <-ticker.C:
			continue

		case <-ctx.Done():
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
			continue
		}
	}
}

func (csp *ConduitStreamProcessor) GenerateSocketListeningRoutines(ctx context.Context) {
	for i := 0; i < csp.session.GetConnectionCount(); i++ {
		go csp.ListenAndHandleDataChannel(ctx, i)
	}
}
