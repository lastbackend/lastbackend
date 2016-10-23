package model

import (
	"time"
)

type ImageTemplate struct {
	UUID        string
	Hub         string
	Owner       string
	Repo        string
	Tag         string
	Memory      int64
	CMD         string
	Entrypoint  string
	Description string
	Category    *ImageCategory
	Env         []string `json:"environments"`
	Ports       []Port   `json:"ports"`
	Volumes     []Volume `json:"volumes"`
	Created     time.Time
	Updated     time.Time
}

type ImageTemplates []ImageTemplate
