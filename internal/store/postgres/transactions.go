package postgres

import (
	"fmt"

	"github.com/volatrade/conduit/internal/models"
)

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
		postgres.kstats.Increment(fmt.Sprintf("conduit.duplicate_inserts.%s", transaction.Pair), 1.0)
	}

	return err

}

func (postgres *DB) BulkInsertTransactions(transactionList []*models.Transaction) error {

	tx := postgres.DB.MustBegin()
	stmt, err := tx.PrepareNamed(TRANSACTION_INSERTION_QUERY)
	if err != nil {
		return err
	}
	defer stmt.Close()
	println(len(transactionList))
	postgres.logger.Infow("bulk insert", "type", "transaction", "size", len(transactionList))
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
			postgres.kstats.Increment(fmt.Sprintf("conduit.duplicate_inserts.tx.%s", transaction.Pair), 1.0)
		}
		if err != nil {
			tx.Rollback()
			return err
		}

	}

	return tx.Commit()

}
