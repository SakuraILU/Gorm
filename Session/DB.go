package session

import "database/sql"

type DB interface {
	Exec(sql string, values ...any) (sql.Result, error)
	Query(sql string, values ...any) (*sql.Rows, error)
	QueryRow(sql string, values ...any) *sql.Row
}
