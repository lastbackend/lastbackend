package model

import (
	"time"
)

type Usage struct {
	UUID      string
	UserID    string
	ServiceID string
	Unit      int64
	Cost      float64
	Price     float64
	Start     time.Time
	Stop      time.Time
	Created   time.Time
	Updated   time.Time
}

type Usages []Usage
