package sqlf

import (
	"database/sql"
)

func (s SQL) Exec(db *sql.DB) (sql.Result, error) {
	query, args := s.BuildSQL()
	return db.Exec(query, args...)
}

func (s SQL) Query(db *sql.DB) (*sql.Rows, error) {
	query, args := s.BuildSQL()
	return db.Query(query, args...)
}

func (s SQL) QueryRow(db *sql.DB) *sql.Row {
	query, args := s.BuildSQL()
	return db.QueryRow(query, args...)
}
