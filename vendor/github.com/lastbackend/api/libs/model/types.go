package model

import (
	"database/sql"
	"encoding/json"
)

// This null types are made for proper json marshalling
type NullString struct {
	sql.NullString
}

func (r NullString) MarshalJSON() ([]byte, error) {
	if !r.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(r.String)
}

type NullInt64 struct {
	sql.NullInt64
}

func (r NullInt64) MarshalJSON() ([]byte, error) {
	if !r.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(r.Int64)
}

type NullBool struct {
	sql.NullBool
}

func (r NullBool) MarshalJSON() ([]byte, error) {
	if !r.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(r.Bool)
}

type NullFloat64 struct {
	sql.NullFloat64
}

func (r NullFloat64) MarshalJSON() ([]byte, error) {
	if !r.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(r.Float64)
}
