package model

import (
	"time"
)

type Stack struct {
	UUID      string
	UserID    string
	Name      string
	Stackfile string
	Services  *Services
	Created   time.Time
	Updated   time.Time
}

type Stacks []Stack
