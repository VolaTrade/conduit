package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/volatrade/conduit/internal/models"
)

const (
	CANDLESTICK_INSERTION_QUERY = `INSERT INTO candlesticks (open, high, low, close, volume, timestamp, pair) VALUES (:open, :high, :low, :close, :volume, :timestamp, :pair) ON CONFLICT DO NOTHING;`
)

func (postgres *DB) BulkInsertCandlestickRows(candlestickRows []*models.Kline) error {

	tx := postgres.DB.MustBegin()

	rollBack := func(tx *sqlx.Tx) {

		if err := tx.Rollback(); err != nil {
			postgres.logger.Errorw("Could not rollback pg transaction", "error", err.Error())
		}
	}

	stmt, err := tx.PrepareNamed(CANDLESTICK_INSERTION_QUERY)

	if err != nil {
		return err
	}

	defer stmt.Close()

	postgres.logger.Infow("bulk insert", "type", "candlestick", "size", len(candlestickRows))

	for _, candlestickRow := range candlestickRows {

		if candlestickRow == nil {
			continue
		}
		result, err := stmt.Exec(candlestickRow)

		if err != nil {
			postgres.logger.Errorw(err.Error(), "description", "executing query")
			rollBack(tx)
			return err
		}

		if rows, err := result.RowsAffected(); rows == 0 && err == nil {
			postgres.kstats.Increment(fmt.Sprintf("duplicate_inserts.cd.%s", candlestickRow.Pair), 1.0)
		}
		if err != nil {
			rollBack(tx)
			return err
		}

	}

	return tx.Commit()

}
