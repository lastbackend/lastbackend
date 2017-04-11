package v1

import (
	"time"
)

type Namespace struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

type NamespaceList []Namespace
