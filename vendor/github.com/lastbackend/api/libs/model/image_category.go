package model

import (
	"time"
)

type ImageCategory struct {
	UUID        string
	Name        string
	Description string
	Created     time.Time
	Updated     time.Time
}

type ImageCategories []ImageCategory
