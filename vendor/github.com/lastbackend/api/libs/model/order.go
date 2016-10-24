package model

import (
	"time"
)

type Order struct {
	UUID          string
	UserID        string
	PaymentSystem string
	Amount        float64
	OrderID       string
	PaymentType   string
	SkuID         string
	Created       time.Time
}

type Orders []Order
