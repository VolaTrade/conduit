package postgres

import (
	"fmt"

	"github.com/volatrade/conduit/internal/models"
)

const (
	TRANSACTION_INSERTION_QUERY = `INSERT INTO transactions(trade_id, time_stamp, pair, price, quantity, is_maker) VALUES(:id, :timestamp, :pair, :price, :quant, :maker) ON CONFLICT DO NOTHING`
	DEPTH_INSERTION_QUERY       = `INSERT INTO order_books (id, pair, time_stamp, bids, asks) VALUES (:id, :pair, :timestamp, :bids, :asks) ON CONFLICT DO NOTHING;`
)

func (postgres *DB) InsertOrderBookRow(obUpdate *models.OrderBookRow) error {
	stmt, err := postgres.DB.PrepareNamed(DEPTH_INSERTION_QUERY)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(obUpdate)
	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); rows == 0 && err == nil {
		postgres.kstats.Increment(fmt.Sprintf("conduit.duplicate_inserts.%s", obUpdate.Pair), 1.0)
	}

	return nil
}

func (postgres *DB) BulkInsertOrderBookRows(orderBookRows []*models.OrderBookRow) error {

	tx := postgres.DB.MustBegin()

	stmt, err := tx.PrepareNamed(DEPTH_INSERTION_QUERY)

	if err != nil {
		return err
	}

	defer stmt.Close()

	postgres.logger.Infow("bulk insert", "type", "orderbook", "size", len(orderBookRows))

	for _, orderBookRow := range orderBookRows {

		if orderBookRow == nil {
			continue
		}

		result, err := stmt.Exec(orderBookRow)

		if err != nil {
			postgres.logger.Errorw(err.Error(), "description", "executing query")
			tx.Rollback()
			return err
		}

		if rows, err := result.RowsAffected(); rows == 0 && err == nil {
			postgres.kstats.Increment(fmt.Sprintf("conduit.duplicate_inserts.ob.%s", orderBookRow.Pair), 1.0)
		}
		if err != nil {
			tx.Rollback()
			return err
		}

	}

	return tx.Commit()

}
