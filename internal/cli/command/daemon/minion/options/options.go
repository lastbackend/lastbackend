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

package options

import (
	"bytes"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type MinionFlags struct {
	Security struct {
		Token string `yaml:"token" mapstructure:"token"`
	} `yaml:"security" mapstructure:"security"`

	Server struct {
		Host string `yaml:"host" mapstructure:"host"`
		Port uint   `yaml:"port" mapstructure:"port"`
		TLS  struct {
			Verify   bool   `yaml:"verify" mapstructure:"verify"`
			FileCA   string `yaml:"ca" mapstructure:"ca"`
			FileCert string `yaml:"cert" mapstructure:"cert"`
			FileKey  string `yaml:"key" mapstructure:"key"`
		} `yaml:"tls" mapstructure:"tls"`
	} `yaml:"server" mapstructure:"server"`

	Registry struct {
		Config string `yaml:"config" mapstructure:"config"`
	} `yaml:"registry" mapstructure:"registry"`

	Workdir        string `yaml:"workdir" mapstructure:"workdir"`
	ManifestDir    string `yaml:"manifestdir" mapstructure:"manifestdir"`
	DisableSelinux bool   `yaml:"disable-selinux" mapstructure:"disable-selinux"`
	Rootless       bool   `yaml:"rootless" mapstructure:"rootless"`
}

func (cfg *MinionFlags) LoadViper(v *viper.Viper) *viper.Viper {
	v.SetConfigType("yaml")

	buf, err := yaml.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(buf)
	if err := v.ReadConfig(reader); err != nil {
		panic(err)
	}

	return v
}

func NewMinionFlags() *MinionFlags {
	s := new(MinionFlags)
	return s
}

func (f *MinionFlags) AddFlags(mainfs *pflag.FlagSet) {

	fs := pflag.NewFlagSet("", pflag.ExitOnError)

	defer func() {
		fs.VisitAll(func(f *pflag.Flag) {
			if len(f.Deprecated) > 0 {
				f.Hidden = false
			}
		})
		mainfs.AddFlagSet(fs)
	}()

	fs.StringVarP(&f.Security.Token, "access-token", "", "", "Set access token for API server")
	fs.StringVarP(&f.Server.Host, "bind-address", "", "0.0.0.0", "Set bind address for API server")
	fs.UintVarP(&f.Server.Port, "bind-port", "", 2992, "Set listening port binding for API server")
	fs.BoolVarP(&f.Server.TLS.Verify, "api-tls-verify", "", false, "Enable check tls for API server")
	fs.StringVarP(&f.Server.TLS.FileCA, "api-tls-ca-file", "", "", "Set path to ca file for API server")
	fs.StringVarP(&f.Server.TLS.FileCert, "api-tls-private-cert-file", "", "", "Set path to cert file for API server")
	fs.StringVarP(&f.Server.TLS.FileKey, "api-tls-private-key-file", "", "", "Set path to key file for API server")
	fs.StringVarP(&f.Workdir, "workdir", "", "/${HOME}/.lastbackend/", "Set directory path to hold state")
	fs.StringVarP(&f.ManifestDir, "manifestdir", "", "", "Set directory path to manifest")
	fs.StringVarP(&f.Registry.Config, "registry-config-path", "", "", "Registry configuration file path")
	fs.BoolVarP(&f.Rootless, "rootless", "", false, "Run rootless")
	fs.BoolVarP(&f.DisableSelinux, "disable-selinux", "", false, "Disable SELinux in containerd if currently enabled")
}
