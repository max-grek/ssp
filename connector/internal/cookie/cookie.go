package cookie

import (
	"database/sql"
	"test-assignment-cookie-sync/connector/internal/cookie/query"
	"test-assignment-cookie-sync/connector/internal/cookie/tx"
)

type cookieQTX interface {
	query.Querier
	tx.Txr
}

type QTX struct {
	query.QWrapper
	tx.TXWrapper
	conn *sql.DB
}

func New(conn *sql.DB) *QTX {
	var q = query.New(conn)
	return &QTX{
		*q,
		*tx.New(q, conn),
		conn,
	}
}
