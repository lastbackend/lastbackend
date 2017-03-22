package config

import (
"time"
)

// The structure of the config to run the daemon
type Config struct {
	Debug bool `yaml:"debug"`

	TokenSecret string `yaml:"secret"`

	HttpServer struct {
		Port int `yaml:"port"`
	} `yaml:"http_server"`

	Etcd struct {
		Endpoints []string      `yaml:"endpoints"`
		TimeOut   time.Duration `yaml:"timeout"`
	} `yaml:"etcd"`

	Docker struct {
		Endpoint, CA, Cert, Key string
	}
}

