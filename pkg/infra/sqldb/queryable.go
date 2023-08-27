package sqldb

import (
	"database/sql"
)

// QueryAble satisfies the *sql.DB and *sql.Tx object methods,
// allowing users to write functions that can take either the DB or the Tx
type QueryAble interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
