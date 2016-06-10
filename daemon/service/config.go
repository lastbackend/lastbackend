package service

import (
	"github.com/deployithq/deployit/daemon/env"
	"fmt"
)

type Config struct {
	Env     []string  `json:"env" yaml:"env"`
	Ports   []string  `json:"ports" yaml:"ports"`
	Volumes []string  `json:"volumes" yaml:"volumes"`
	CMD     []string  `json:"cmd" yaml:"cmd"`
	Memory  int64     `json:"memory" yaml:"memory"`
	Image   string    `json:"image" yaml:"image"`
}

var configs map[string]*Config

func init() {
	fmt.Println(`Init service configs`)
	configs = make(map[string]*Config)

	configs[`redis`] = &Config{
		Image: `library/redis`,
		Ports: []string{"6379"},
	}

	configs[`memcached`] = &Config{
		Image: `library/memcached`,
		Ports: []string{"11211"},
	}

	configs[`postgres`] = &Config{
		Image: `library/postgres`,
		Ports: []string{"5432"},
	}

	configs[`couchdb`] = &Config{
		Image: `library/couchdb`,
		Ports: []string{"5984"},
	}

	configs[`mysql`] = &Config{
		Image: `library/mysql`,
		Ports: []string{"3306"},
	}

	configs[`mongo`] = &Config{
		Image: `library/mongo`,
		Ports: []string{"27017"},
	}
}

// Todo: select config for services
func (c *Config) Get(e *env.Env, name string) error {
	e.Log.Info(`Get config for `, name)

	if val, ok := configs[name]; ok {
		*c = *val
	}

	return nil
}
