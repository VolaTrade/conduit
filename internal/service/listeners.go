package service

import (
	"os"
	"sync"
)

func (ts *ConduitService) ListenAndHandleDataChannels(index int, wg *sync.WaitGroup, quit chan bool) {
	defer wg.Done()
	for {
		select {

		case <-quit:
			return

		case transaction := <-ts.transactionChannels[index]:
			ts.handleTransaction(transaction, index)

		case orderBookRow := <-ts.orderBookChannels[index]:
			ts.handleOrderBookRow(orderBookRow, index)
		}

	}
}

//TODO there's a better way to structure this
func (ts *ConduitService) ListenForDatabasePriveleges(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	var err error
	for {
		if _, writeToCache := os.Stat("start"); writeToCache == nil {
			ts.logger.Infow("establishing database connections, moving cache to databse, and purging cache")
			ts.dbStreams.MakeConnections()
			ts.writeToDB = true
			cachedTransactions := ts.cache.GetAllTransactions()
			cachedOrderBookRows := ts.cache.GetAllOrderBookRows()
			err = ts.dbStreams.TransferTransactionCache(cachedTransactions)
			if err != nil {
				panic(err)
			}

			err = ts.dbStreams.TransferOrderBookCache(cachedOrderBookRows)

			if err != nil {
				panic(err)
			}

			ts.cache.Purge()
			return
		}
	}
}

func (ts *ConduitService) ListenForExit(wg *sync.WaitGroup, exit func()) {
	defer wg.Done()
	for {
		if _, err := os.Stat("finish"); err == nil {
			ts.logger.Infow("Finish signal recieved")
			exit()
			return
		}
	}
}
