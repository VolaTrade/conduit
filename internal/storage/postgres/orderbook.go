package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/volatrade/conduit/internal/models"
)

const (
	DEPTH_INSERTION_QUERY = `INSERT INTO order_books (id, pair, time_stamp, bids, asks) VALUES (:id, :pair, :timestamp, :bids, :asks) ON CONFLICT DO NOTHING;`
)

func (postgres *DB) BulkInsertOrderBookRows(orderBookRows []*models.OrderBookRow) error {

	tx := postgres.DB.MustBegin()

	rollBack := func(tx *sqlx.Tx) {

		if err := tx.Rollback(); err != nil {
			postgres.logger.Errorw("Could not rollback pg transaction", "error", err.Error())
		}
	}

	stmt, err := tx.PrepareNamed(DEPTH_INSERTION_QUERY)

	if err != nil {
		return err
	}

	defer stmt.Close()

	postgres.logger.Infow("bulk insert", "type", "orderbok", "size", len(orderBookRows))

	for _, orderBookRow := range orderBookRows {

		if orderBookRow == nil {
			continue
		}
		result, err := stmt.Exec(orderBookRow)

		if err != nil {
			postgres.logger.Errorw(err.Error(), "description", "executing query")
			rollBack(tx)
			return err
		}

		if rows, err := result.RowsAffected(); rows == 0 && err == nil {
			postgres.kstats.Increment(fmt.Sprintf("duplicate_inserts.ob.%s", orderBookRow.Pair), 1.0)
		}
		if err != nil {
			rollBack(tx)
			return err
		}

	}

	return tx.Commit()

}
