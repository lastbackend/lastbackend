package app

import (
	"fmt"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/utils"
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

func (c *Config) Create(e *env.Env, hub, name, tag string) error {
	e.Log.Info(`Sync config`)

	if hub == `` {
		hub = env.Default_hub
	}

	c.Image = fmt.Sprintf("%s/%s:%s", hub, name, tag)
	c.Created = time.Now()

	return nil
}

func (c *Config) Sync(e *env.Env, layer string) error {
	e.Log.Info(`Sync config`)

	path := fmt.Sprintf("%s/apps/%s", env.Default_root_path, layer)

	config := new(Config)
	if err := utils.ReadConfig(path, &config); err != nil {
		return err
	}

	c.Env = config.Env
	c.Ports = config.Ports
	c.Volumes = config.Volumes
	c.CMD = config.CMD
	c.Memory = config.Memory
	c.Updated = time.Now()

	return nil
}
