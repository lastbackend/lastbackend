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
	DefaultBindServerPort    = 2992
	DefaultWorkDir           = "/${HOME}/.lastbackend/"
	DefaultCIDR              = "172.0.0.0/24"
)

type Config struct {
	Debug bool `yaml:"debug,omitempty"`

	Security struct {
		Token string `yaml:"token,omitempty"`
	} `yaml:"security,omitempty"`

	Server struct {
		Host string `yaml:"host,omitempty"`
		Port uint   `yaml:"port,omitempty"`
		TLS  struct {
			Verify   bool   `yaml:"verify,omitempty"`
			FileCA   string `yaml:"ca,omitempty"`
			FileCert string `yaml:"cert,omitempty"`
			FileKey  string `yaml:"key,omitempty"`
		} `yaml:"tls,omitempty"`
	} `yaml:"server,omitempty"`

	Registry struct {
		Config string `yaml:"config,omitempty"`
	} `yaml:"registry,omitempty"`

	WorkDir        string `yaml:"workdir,omitempty"`
	ManifestDir    string `yaml:"manifestdir,omitempty"`
	CIDR           string `yaml:"cidr,omitempty"`
	DisableSeLinux bool   `yaml:"disable-selinux,omitempty"`
	Rootless       bool   `yaml:"rootless,omitempty"`
}
