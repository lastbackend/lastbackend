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

const (
	DefaultBindServerAddress = "0.0.0.0"
	DefaultBindServerPort    = 2967
	DefaultInternalDomain    = "lb.local"
	DefaultWorkDir           = "/${HOME}/.lastbackend/"
)

type Config struct {
	Debug bool `yaml:"debug,omitempty"`

	Security struct {
		Token string `yaml:"token,omitempty"`
	} `yaml:"security,omitempty"`

	WorkDir            string `yaml:"workdir"`
	ClusterName        string `yaml:"name"`
	ClusterDescription string `yaml:"description"`
	Rootless           bool   `yaml:"rootless"`

	Server ServerConfig `yaml:"server,omitempty"`
	Vault  VaultConfig  `yaml:"vault"`
	Domain DomainConfig `yaml:"domain"`
}

type ServerConfig struct {
	Host string `yaml:"host,omitempty"`
	Port uint   `yaml:"port,omitempty"`
	TLS  struct {
		Verify   bool   `yaml:"verify,omitempty"`
		FileCA   string `yaml:"ca,omitempty"`
		FileCert string `yaml:"cert,omitempty"`
		FileKey  string `yaml:"key,omitempty"`
	} `yaml:"tls,omitempty"`
}

type VaultConfig struct {
	Token    string `yaml:"token"`
	Endpoint string `yaml:"endpoint"`
}

type DomainConfig struct {
	Internal string `yaml:"internal"`
	External string `yaml:"external"`
}
