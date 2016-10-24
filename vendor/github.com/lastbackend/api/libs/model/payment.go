package model

import (
	"time"
)

type Payment struct {
	UUID          string
	UserID        string
	PaymentSystem string
	Payload       string
	Created       time.Time
}

type Payments []Payment
