package model

import "time"

type Node struct {
	UUID     string    `json:"uuid,omitempty"`
	UserID   string    `json:"user,omitempty"`
	Provider *Provider `json:"provider,omitempty"`
	Token    string    `json:"token,omitempty"`
	RegionID string    `json:"region,omitempty"`
	Version  string    `json:"version,omitempty"`
	Online   bool      `json:"online,omitempty"`
	IPs      []struct {
		ID        string `json:"id,omitempty"`
		Interface string `json:"interface,omitempty"`
		IP        string `json:"ip,omitempty"`
		Main      bool   `json:"main,omitempty"`
	} `json:"ips,omitempty"`

	Hostname string `json:"hostname,omitempty"`
	OS       struct {
		Name string `json:"name,omitempty"`
		Type string `json:"type,omitempty"`
	} `json:"os,omitempty"`

	IP struct {
		External string `json:"external,omitempty"`
		Internal string `json:"internal,omitempty"`
		Local    string `json:"local,omitempty"`
	} `json:"ip,omitempty"`

	Architecture string `json:"architecture,omitempty"`

	Driver struct {
		Name    string `json:"name,omitempty"`
		Version string `json:"version,omitempty"`
	} `json:"driver,omitempty"`

	Storage struct {
		Available string `json:"available,omitempty"`
		Used      string `json:"used,omitempty"`
		Total     string `json:"total,omitempty"`
	} `json:"storage,omitempty"`

	CPU struct {
		Name  string `json:"name,omitempty"`
		Cores int64  `json:"cores,omitempty"`
	} `json:"cpu,omitempty"`

	Memory struct {
		Total int64 `json:"total,omitempty"`
		Used  int64 `json:"used,omitempty"`
	} `json:"memory,omitempty"`

	Containers struct {
		Started int64 `json:"started,omitempty"`
		Total   int64 `json:"total,omitempty"`
		Errored int64 `json:"errored,omitempty"`
	} `json:"containers,omitempty"`

	Deleted     bool      `json:"deleted,omitempty"`
	Created     time.Time `json:"created,omitempty"`
	Updated     time.Time `json:"updated,omitempty"`
	Description string    `json:"description,omitempty"`

	Docker struct {
		Host string `json:"-"`
		Cert string `json:"-"`
		Key  string `json:"-"`
		CA   string `json:"-"`
		TLS  bool   `json:"-"`
	} `json:"-"`
}

type Nodes []Node
