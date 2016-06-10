package service

import (
	"fmt"
	"github.com/deployithq/deployit/daemon/env"
	"time"
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

func (c *Config) Create(e *env.Env, name string) error {
	e.Log.Info(`Sync config`)

	c.Image = fmt.Sprintf("%s/%s:%s", `library`, name, `latest`)
	c.Created = time.Now()

	return nil
}
