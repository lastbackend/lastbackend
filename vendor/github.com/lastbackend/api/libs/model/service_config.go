package model

type ServiceConfig struct {
	UUID        string   `json:"uuid,omitempty"`
	ServiceID   string   `json:"service,omitempty"`
	ClusterID   string   `json:"cluster,omitempty"`
	CloudRegion string   `json:"cloud_region,omitempty"`
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
	UseCloud    bool     `json:"use_cloud,omitempty"`
}

type Port struct {
	Host      int    `json:"host"`
	Container int    `json:"container"`
	Protocol  string `json:"protocol"`
}

type Ports []Port

type Volume struct {
	Host      string `json:"host"`
	Container string `json:"container"`
}

type Volumes []Volume

type ServiceConfigs []ServiceConfig
