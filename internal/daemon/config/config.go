//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package config

import (
	agent_config "github.com/lastbackend/lastbackend/internal/agent/config"
	server_config "github.com/lastbackend/lastbackend/internal/server/config"
)

const (
	DefaultBindServerAddress = "0.0.0.0"
	DefaultBindServerPort    = 2992
	DefaultCIDR              = "172.0.0.0/24"
	DefaultInternalDomain    = "lb.local"
	DefaultRootDir           = "/var/lib/lastbackend/"
)

type Config struct {
	Debug           bool   `yaml:"debug"`
	DisableSchedule bool   `yaml:"no-schedule"`
	DisableServer   bool   `yaml:"agent"`
	RootDir         string `yaml:"root-dir"`

	StorageDriver string `yaml:"storage-driver"`
	ManifestDir   string `yaml:"manifest-dir"`
	CIDR          string `yaml:"cidr"`

	Security   Security     `yaml:"security"`
	APIServer  ServerConfig `yaml:"api-server"`
	NodeServer ServerConfig `yaml:"node-server"`
	NodeClient NodeClient   `yaml:"node-client"`
	Vault      VaultConfig  `yaml:"vault"`
	Domain     DomainConfig `yaml:"domain"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port uint   `yaml:"port"`
	TLS  struct {
		Verify   bool   `yaml:"verify"`
		FileCA   string `yaml:"ca"`
		FileCert string `yaml:"cert"`
		FileKey  string `yaml:"key"`
	} `yaml:"tls"`
}

type VaultConfig struct {
	Token    string `yaml:"token"`
	Endpoint string `yaml:"endpoint"`
}

type DomainConfig struct {
	Internal string `yaml:"internal"`
	External string `yaml:"external"`
}

type Security struct {
	Token string `yaml:"token"`
}

type NodeClient struct {
	Address string `yaml:"uri"`
	TLS     struct {
		Verify   bool   `yaml:"verify"`
		FileCA   string `yaml:"ca"`
		FileCert string `yaml:"cert"`
		FileKey  string `yaml:"key"`
	} `yaml:"tls"`
}

func (c Config) GetServerConfig() server_config.Config {
	cfg := server_config.Config{}
	cfg.Debug = c.Debug
	cfg.RootDir = c.RootDir
	cfg.Security = server_config.SecurityConfig(c.Security)
	cfg.Server = server_config.ServerConfig(c.APIServer)
	cfg.Vault = server_config.VaultConfig(c.Vault)
	cfg.Domain = server_config.DomainConfig(c.Domain)
	return cfg
}

func (c Config) GetAgentConfig() agent_config.Config {
	cfg := agent_config.Config{}
	cfg.Debug = c.Debug
	cfg.RootDir = c.RootDir
	cfg.StorageDriver = c.StorageDriver
	cfg.ManifestDir = c.ManifestDir
	cfg.CIDR = c.CIDR
	cfg.Security = agent_config.SecurityConfig(c.Security)
	cfg.Server = agent_config.ServerConfig(c.NodeServer)
	cfg.API = agent_config.NodeClient(c.NodeClient)
	return cfg
}
