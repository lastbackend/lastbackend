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

	Config     Config     `json:"config,omitempty"`
	HostConfig HostConfig `json:"host_config,omitempty"`

	State State  `json:"state,omitempty"`
	Ports []Port `json:"ports,omitempty"`
}

type Port struct {
	Container int64  `json:"pivate,omitempty" yaml:"pivate,omitempty"`
	Host      int64  `json:"public,omitempty"  yaml:"public,omitempty"`
	Protocol  string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}

type Config struct {
	Image      string   `json:"image" yaml:"image,omitempty"`
	Env        []string `json:"env" yaml:"env,omitempty"`
	Cmd        []string `json:"cmd" yaml:"cmd,omitempty"`
	Volumes    []string `json:"volumes" yaml:"volumes,omitempty"` // []string{"/data:/data:rw"}
	Ports      []string `json:"ports" yaml:"ports,omitempty"` // []string{"80:80"}
	Memory     int64    `json:"memory" yaml:"memory,omitempty"`
	Entrypoint []string `json:"entrypoint" yaml:"entrypoint,omitempty"`
}

type HostConfig struct {
	Binds         []string            `json:"binds" yaml:"binds,omitempty"`
	Ports         []string            `json:"ports" yaml:"ports,omitempty"` // []string{"80:80"}
	RestartPolicy RestartPolicyConfig `json:"restart" yaml:"restart,omitempty"`
	Memory        int64               `json:"memory" yaml:"memory,omitempty"`
	Privileged    bool                `json:"privileged" yaml:"privileged,omitempty"`
}

type Volume struct {
	Host      string `json:"host" yaml:"host,omitempty"`
	Container string `json:"container" yaml:"container,omitempty"`
	Options   string `json:"options" yaml:"options,omitempty"`
}

type RestartPolicyConfig struct {
	Name    string `json:"name" yaml:"name,omitempty"`
	Attempt int    `json:"attempt" yaml:"attempt,omitempty"`
}

type State struct {
	Running    bool      `json:"running,omitempty" yaml:"running,omitempty"`
	Paused     bool      `json:"paused,omitempty" yaml:"paused,omitempty"`
	Restarting bool      `json:"restarting,omitempty" yaml:"restarting,omitempty"`
	OOMKilled  bool      `json:"oomKilled,omitempty" yaml:"oomKilled,omitempty"`
	Pid        int       `json:"pid,omitempty" yaml:"pid,omitempty"`
	ExitCode   int       `json:"exit_code,omitempty" yaml:"exit_code,omitempty"`
	Error      string    `json:"error,omitempty" yaml:"error,omitempty"`
	Started    time.Time `json:"started,omitempty" yaml:"started,omitempty"`
	Finished   time.Time `json:"finished,omitempty" yaml:"finished,omitempty"`
}

type Image struct {
	Name string     `json:"name" yaml:"name,omitempty"`
	Auth AuthConfig `json:"auth,omitempty" yaml:"auth,omitempty"`
}

type AuthConfig struct {
	Username string `json:"username" yaml:"username,omitempty"`
	Password string `json:"password" yaml:"password,omitempty"`
	Email    string `json:"email" yaml:"email,omitempty"`
	Host     string `json:"serveraddres" yaml:"serveraddres,omitempty"`
}

type BuildImageOptions struct {
	Name           string    `json:"name"`
	RmTmpContainer bool      `json:"rm"`
	ContextDir     string    `json:"context"`
	RawJSONStream  bool      `json:"raw"`
	InputStream    io.Reader `json:"-"`
	OutputStream   io.Writer `json:"-"`
}
