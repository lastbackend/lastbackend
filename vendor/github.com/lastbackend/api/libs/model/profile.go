package model

import (
	"time"
)

type UserProfile struct {
	UUID      string
	UserID    string
	FirstName string
	LastName  string
	Country   string
	City      string
	State     string
	ZipCode   string
	Address   string
	Phone     string
	Company   string
	Created   time.Time
	Updated   time.Time
}
