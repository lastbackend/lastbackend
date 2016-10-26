package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type nullString struct {
	sql.NullString
}

func (ns nullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns nullString) UnmarshalJSON(v interface{}) error {
	if !ns.Valid {
		return nil
	}
	return json.Unmarshal([]byte(ns.String), v)
}

type nullInt64 struct {
	sql.NullInt64
}

type nullBool struct {
	sql.NullBool
}

type nullFloat64 struct {
	sql.NullFloat64
}

type nullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *nullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt nullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
