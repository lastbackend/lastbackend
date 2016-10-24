package model

import (
	"time"
)

type Service struct {
	UUID        string           `json:"uuid,omitempty"`
	NodeID      string           `json:"node,omitempty"`
	ClusterID   string           `json:"cluster,omitempty"`
	Name        string           `json:"name,omitempty"`
	Description string           `json:"description,omitempty"`
	Domain      string           `json:"domain,omitempty"`
	Status      string           `json:"status,omitempty"`
	Message     string           `json:"message,omitempty"`
	Pull        bool             `json:"pull,omitempty"`
	Provision   bool             `json:"provider,omitempty"`
	Deleted     bool             `json:"deleted,omitempty"`
	Started     time.Time        `json:"started,omitempty"`
	Created     time.Time        `json:"created,omitempty"`
	Updated     time.Time        `json:"updated,omitempty"`
	User        *User            `json:"user,omitempty"`
	Image       *Image           `json:"image,omitempty"`
	Domains     *[]ServiceDomain `json:"domains,omitempty"`
	Config      *ServiceConfig   `json:"config,omitempty"`
	Source      *ServiceSource   `json:"source,omitempty"`
	Containers  struct {
		Total   int64 `json:"total,omitempty"`
		Started int64 `json:"started,omitempty"`
		Errored int64 `json:"errored,omitempty"`
	} `json:"containers,omitempty"`
}

type Services []Service

type ServiceByNode struct {
	ID   string
	Name NullString
}

type ServicesByNode []ServiceByNode
