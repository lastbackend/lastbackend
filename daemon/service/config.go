package service

import (
	"time"
	"github.com/deployithq/deployit/daemon/env"
)

type Config struct {
	Env     []string  `json:"env" yaml:"env"`
	Ports   []string  `json:"ports" yaml:"ports"`
	Volumes []string  `json:"volumes" yaml:"volumes"`
	CMD     []string  `json:"cmd" yaml:"cmd"`
	Memory  int64     `json:"memory" yaml:"memory"`
	Image   string    `json:"image" yaml:"image"`
	Created time.Time `json:"created" yaml:"created"`
	Updated time.Time `json:"updated" yaml:"updated"`
}

func (c *Config) Get(e *env.Env, name string) error {
	return nil
}