package model

type Config struct {
	Replicas   int32          `json:"replicas,omitempty"`
	Command    []string       `json:"command,omitempty"`
	Args       []string       `json:"args,omitempty"`
	WorkingDir string         `json:"workdir,omitempty"`
	Ports      []PortConfig   `json:"ports,omitempty"`
	Env        []EnvVarConfig `json:"env,omitempty"`
	Volumes    []VolumeConfig `json:"volumes,omitempty"`
}

type PortConfig struct {
	Name          string `json:"name,omitempty"`
	HostPort      int32  `json:"host,omitempty"`
	ContainerPort int32  `json:"container"`
	Protocol      string `json:"protocol,omitempty"`
	HostIP        string `json:"ip,omitempty"`
}

type EnvVarConfig struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

type VolumeConfig struct {
	Name      string `json:"name"`
	ReadOnly  bool   `json:"readonly,omitempty"`
	MountPath string `json:"mountpath"`
	SubPath   string `json:"subpath,omitempty"`
}
