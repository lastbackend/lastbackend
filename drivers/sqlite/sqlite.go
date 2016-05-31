package sqlite

import (
	"database/sql"
	"github.com/deployithq/deployit/drivers/interfaces"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	conn *sql.DB
}

func Open(dbPath string) (*SQLite, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return new(SQLite), err
	}

	sqlite := SQLite{}
	sqlite.conn = conn

	return &sqlite, nil
}

func (db *SQLite) Exec(query string, args ...interface{}) (interfaces.Result, error) {
	result, err := db.conn.Exec(query, args...)
	if err != nil {
		return new(SQLiteResult), err
	}

	sqlRow := SQLiteResult{result}

	return sqlRow, nil
}

func (db *SQLite) Query(query string, args ...interface{}) (interfaces.Rows, error) {
	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return new(SQLiteRows), err
	}

	sqlRows := new(SQLiteRows)
	sqlRows.Rows = rows
	return sqlRows, nil
}

func (db *SQLite) QueryRow(query string, args ...interface{}) interfaces.Row {
	row := db.conn.QueryRow(query, args...)
	return &SQLiteRow{row}
}

type SQLiteRows struct {
	Rows *sql.Rows
}

func (r SQLiteRows) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

func (r SQLiteRows) Next() bool {
	return r.Rows.Next()
}

func (r SQLiteRows) Err() error {
	return r.Rows.Err()
}

func (r SQLiteRows) Columns() ([]string, error) {
	return r.Rows.Columns()
}

func (r SQLiteRows) Close() error {
	return r.Rows.Close()
}

type SQLiteRow struct {
	Row *sql.Row
}

func (r *SQLiteRow) Scan(dest ...interface{}) error {
	return r.Row.Scan(dest...)
}

type SQLiteResult struct {
	sql.Result
}

func (r SQLiteResult) LastInsertId() (int64, error) {
	return r.LastInsertId()
}

func (r SQLiteResult) RowsAffected() (int64, error) {
	return r.RowsAffected()
}
