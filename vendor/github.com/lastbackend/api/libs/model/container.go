package model

import (
	"time"
)

type Container struct {
	UUID      string `json:"uuid,omitempty"`
	UserID    string `json:"user,omitempty"`
	ServiceID string `json:"service,omitempty"`
	NodeID    string `json:"node,omitempty"`
	CID       string `json:"cid,omitempty"`
	Image     struct {
		Hub   string `json:"hub,omitempty"`
		Owner string `json:"owner,omitempty"`
		Repo  string `json:"repo,omitempty"`
		Tag   string `json:"tag,omitempty"`
	} `json:"image,omitempty"`
	Error      string    `json:"error,omitempty"`
	Status     string    `json:"status,omitempty"`
	PID        int64     `json:"pid,omitempty"`
	ExitCode   int64     `json:"exit_code,omitempty"`
	Running    bool      `json:"running,omitempty"`
	Paused     bool      `json:"paused,omitempty"`
	Restarting bool      `json:"restarting,omitempty"`
	OOMKilled  bool      `json:"oomkilled,omitempty"`
	Name       string    `json:"name,omitempty"`
	Message    string    `json:"message,omitempty"`
	Ports      []Port    `json:"ports,omitempty"`
	Memory     int64     `json:"memory,omitempty"`
	Deleted    bool      `json:"deleted,omitempty"`
	Started    time.Time `json:"started,omitempty"`
	Finished   time.Time `json:"finished,omitempty"`
	Created    time.Time `json:"created,omitempty"`
	Updated    time.Time `json:"updated,omitempty"`
	Config     struct {
		CMD         string   `json:"cmd"`
		Network     string   `json:"network"`
		Strategy    string   `json:"strategy"`
		PID         string   `json:"pid"`
		EntryPoint  string   `json:"entrypoint"`
		Ports       []Port   `json:"ports"`
		Volumes     []Volume `json:"volumes"`
		Env         []string `json:"env"`
		Scale       int64    `json:"scale"`
		Memory      int64    `json:"memory"`
		AutoDeploy  bool     `json:"autodeploy"`
		AutoRestart bool     `json:"autorestart"`
		AutoDestroy bool     `json:"autodestroy"`
		Privileged  bool     `json:"privileged"`
	} `json:"config"`
}

type Containers []Container
