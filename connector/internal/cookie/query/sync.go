package query

import (
	"context"
	"database/sql"
)

type Querier interface {
	PersistNotify(ctx context.Context, id string, cid string) error
	PersistCookie(ctx context.Context, id string, cid string) error

	upgrageToTx(tx *sql.Tx)
}

type QueryTxFunc func(Querier) error

type txdb interface {
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type QWrapper struct {
	txdb
}

func New(db txdb) *QWrapper {
	return &QWrapper{db}
}

func (q QWrapper) PersistNotify(ctx context.Context, id, cid string) error {
	query := `INSERT INTO ssp.cookie (id, dsp_cookie_id) VALUES ($1, $2)`
	_, err := q.ExecContext(ctx, query, id, cid)
	return err
}

func (q QWrapper) PersistCookie(ctx context.Context, id, cid string) error {
	query := `INSERT INTO ssp.cookie (id, ssp_cookie_id) VALUES ($1, $2)`
	_, err := q.ExecContext(ctx, query, id, cid)
	return err
}

func (q *QWrapper) upgrageToTx(tx *sql.Tx) {
	q.txdb = tx
}

func UpgradeToTx(tx *sql.Tx, q Querier) {
	q.upgrageToTx(tx)
}
