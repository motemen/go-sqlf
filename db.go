// Convenient shortcut methods in combination with database/sql.

package sqlf

import (
	"database/sql"
)

type dbQueryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func (s SQL) Exec(db dbQueryer) (sql.Result, error) {
	query, args := s.BuildSQL()
	return db.Exec(query, args...)
}

func (s SQL) Query(db dbQueryer) (*sql.Rows, error) {
	query, args := s.BuildSQL()
	return db.Query(query, args...)
}

func (s SQL) QueryRow(db dbQueryer) *sql.Row {
	query, args := s.BuildSQL()
	return db.QueryRow(query, args...)
}
