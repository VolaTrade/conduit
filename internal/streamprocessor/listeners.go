package streamprocessor
import (
	"context"
	"os"
	"sync"
	"time"
)

//ListenAndHandleDataChannels waits for transaction or orderbook to come in from their respective channel before invoking handler function
func (csp *ConduitStreamProcessor) ListenAndHandleDataChannels(ctx context.Context, index int, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	txChannel, obChannel := csp.transactionChannels[index], csp.orderBookChannels[index]
	for {
		select {

		case transaction := <-txChannel:
			csp.handleTransaction(transaction, index)

		case orderBookRow := <-obChannel:
			csp.handleOrderBookRow(orderBookRow, index)

		case <-ctx.Done():
			return
		}

	}
}

//ListenForDatabasePriveleges checks every fifteen seconds for prevalance of touch file, signals database writing when found
func (csp *ConduitStreamProcessor) ListenForDatabasePriveleges(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(time.Duration(time.Second * 15))
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
			if err := csp.dbStreams.TransferTransactionCache(csp.cache.GetAllTransactions()); err != nil {
				csp.logger.Errorw(err.Error(), "attempt", attempts)
				attempts -= 1
				continue
			}
			csp.cache.PurgeTransactions()
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
			return

		}
	}
}

//ListenForExit checks every fifteen seconds for prevalance of finish file
func (csp *ConduitStreamProcessor) ListenForExit(ctx context.Context, wg *sync.WaitGroup, exit func()) {
	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(time.Duration(time.Second * 15))

	for {
		if _, err := os.Stat("finish"); err == nil {
			csp.logger.Infow("Finish signal recieved")
			exit()
			return
		}

		select {
		case <-ticker.C:
			continue

		case <-ctx.Done():
			return
		}
	}
}
