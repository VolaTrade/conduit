package postgres

import (
	"fmt"

	"github.com/volatrade/tickers/internal/models"
)

const (
	TRANSACTION_INSERTION_QUERY = `INSERT INTO transactions(trade_id, time_stamp, pair, price, quantity, is_maker) VALUES(:id, :timestamp, :pair, :price, :quant, :maker) ON CONFLICT DO NOTHING`
	DEPTH_INSERTION_QUERY       = `INSERT INTO order_books (id, pair, time_stamp, bids, asks) VALUES (:id, :pair, :timestamp, :bids, :asks);`
)

func (postgres *DB) InsertOrderBookUpdate(obUpdate *models.OrderBookRow) error {
	stmt, err := postgres.DB.PrepareNamed(DEPTH_INSERTION_QUERY)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(obUpdate)
	if err != nil {
		return err
	}
	return nil
}

//InsertTransaction inserts transaction into database
func (postgres *DB) InsertTransaction(transaction *models.Transaction) error {
	stmt, err := postgres.DB.PrepareNamed(TRANSACTION_INSERTION_QUERY)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(transaction)

	if err != nil {
		return err
	}
	if rows, err := result.RowsAffected(); rows == 0 && err == nil {
		postgres.statz.Client.Increment(fmt.Sprintf("tickers.duplicate_inserts.%s", transaction.Pair))
	}

	return err

}

func (postgres *DB) BulkInsertCache(transactionList []*models.Transaction) error {

	tx := postgres.DB.MustBegin()
	stmt, err := tx.PrepareNamed(TRANSACTION_INSERTION_QUERY)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, transaction := range transactionList {

		if transaction == nil {
			continue
		}
		result, err := stmt.Exec(transaction)

		if err != nil {
			tx.Rollback()
			return err
		}

		if rows, err := result.RowsAffected(); rows == 0 && err == nil {
			postgres.statz.Client.Increment(fmt.Sprintf("tickers.duplicate_inserts.%s", transaction.Pair))
		}
		if err != nil {
			tx.Rollback()
			return err
		}

	}

	return tx.Commit()

}
