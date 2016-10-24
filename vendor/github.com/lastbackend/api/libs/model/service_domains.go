package model

import (
	"time"
)

type ServiceDomain struct {
	UUID      string    `json:"uuid,omitempty"`
	ServiceID string    `json:"service,omitempty"`
	Name      string    `json:"name"`
	Main      bool      `json:"default"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}
