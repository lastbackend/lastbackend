package interfaces

import (
	"io"
	"time"
)

type Container struct {
	CID     string `json:"cid"`
	Name    string `json:"name,omitempty"`
	Image   string `json:"image,omitempty"`
	Command string `json:"command,omitempty"`

	State State  `json:"state,omitempty"`
	Ports []Port `json:"ports,omitempty"`

	Config     Config
	HostConfig HostConfig

	LB struct {
		ID  string
		APP string
	} `json:"-"`
}

type Port struct {
	Container int64  `json:"pivate,omitempty"`
	Host      int64  `json:"public,omitempty"`
	Protocol  string `json:"protocol,omitempty"`
}

type Config struct {
	Image      string       `json:"image"`
	Env        []string     `json:"env"`
	Cmd        []string     `json:"cmd"`
	Volumes    []Volume     `json:"volumes"`
	Mounts     []Mount      `json:"mounts"`
	DNS        DNSConfig    `json:"dns"`
	CPU        CPUConfig    `json:"cpu"`
	Ports      []Port       `json:"ports"`
	Memory     MemoryConfig `json:"memory"`
	Entrypoint []string     `json:"entrypoint"`
}

type HostConfig struct {
	DNS           DNSConfig           `json:"dns"`
	Binds         []string            `json:"binds"`
	CPU           CPUConfig           `json:"cpu"`
	Ports         []Port              `json:"ports"`
	RestartPolicy RestartPolicyConfig `json:"restart"`
	Memory        MemoryConfig        `json:"memory"`
	Privileged    bool                `json:"privileged"`
}

type Volume struct {
	Host      string `json:"host"`
	Container string `json:"container"`
	Options   string `json:"options"`
}

type Mount struct {
	Name        string `json:"name"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Driver      string `json:"driver"`
	Mode        string `json:"mode"`
	RW          bool   `json:"rw"`
}

type MemoryConfig struct {
	Total       int64 `json:"total"`
	Swap        int64 `json:"swap"`
	Swappiness  int64 `json:"swappiness"`
	Reservation int64 `json:"reservation"`
	Kernel      int64 `json:"kernel"`
}

type DNSConfig struct {
	Server  []string `json:"server"`
	Options []string `json:"options"`
	Search  []string `json:"search"`
}

type CPUConfig struct {
	Shares int64  `json:"shares,omitempty"`
	Set    string `json:"set,omitempty"`
	CPUs   string `json:"cpus,omitempty"`
	MEMs   string `json:"mems,omitempty"`
	Quota  int64  `json:"quota,omitempty"`
	Period int64  `json:"period,omitempty"`
}

type RestartPolicyConfig struct {
	Name    string `json:"name"`
	Attempt int    `json:"attempt"`
}

type State struct {
	Running    bool      `json:"running,omitempty"`
	Paused     bool      `json:"paused,omitempty"`
	Restarting bool      `json:"restarting,omitempty"`
	OOMKilled  bool      `json:"oomKilled,omitempty"`
	Pid        int       `json:"pid,omitempty"`
	ExitCode   int       `json:"exit_code,omitempty"`
	Error      string    `json:"error,omitempty"`
	Started    time.Time `json:"started,omitempty"`
	Finished   time.Time `json:"finished,omitempty"`
}

type Image struct {
	Name string     `json:"name"`
	Auth AuthConfig `json:"auth,omitempty"`
}

type AuthConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Host     string `json:"serveraddres"`
}

type Node struct {
	UUID string

	Driver struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"driver"`

	Hostname     string `json:"hostname"`
	Architecture string `json:"arhitecture"`

	OS struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"os"`

	CPU struct {
		Name  string `json:"name"`
		Cores int64  `json:"cores"`
	} `json:"cpu"`

	Memory struct {
		Total int64 `json:"total"`
	} `json:"memory"`

	IPs []struct {
		Interface string `json:"interfaces,omitempty"`
		IP        string `json:"ip,omitempty"`
		Main      bool   `json:"main,omitempty"`
	} `json:"ips"`

	Storage struct {
		Available string `json:"available"`
		Used      string `json:"used"`
		Total     string `json:"total"`
	} `json:"storage"`
}

type BuildImageOptions struct {
	Name           string    `json:"name"`
	RmTmpContainer bool      `json:"rm"`
	ContextDir     string    `json:"context"`
	RawJSONStream  bool      `json:"raw"`
	InputStream    io.Reader `json:"-"`
	OutputStream   io.Writer `json:"-"`
}
