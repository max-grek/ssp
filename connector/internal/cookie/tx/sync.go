package tx

import (
	"context"
	"database/sql"
	"fmt"
	"test-assignment-cookie-sync/connector/internal/cookie/query"
)

type Txr interface {
	TestTx(ctx context.Context, str string) error
}

type TXWrapper struct {
	query.Querier
	conn *sql.DB
}

func New(q *query.QWrapper, conn *sql.DB) *TXWrapper {
	return &TXWrapper{q, conn}
}

func (w TXWrapper) execTx(ctx context.Context, fn query.QueryTxFunc) error {
	tx, err := w.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query.UpgradeToTx(nil, w.Querier)
	if err := fn(w.Querier); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
