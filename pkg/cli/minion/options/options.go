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
	AccessToken string `yaml:"token" mapstructure:"token"`

	Server struct {
		Host string `yaml:"host" mapstructure:"host"`
		Port uint   `yaml:"port" mapstructure:"port"`
		TLS  struct {
			Verify   bool   `yaml:"verify" mapstructure:"verify"`
			FileCert string `yaml:"cert" mapstructure:"cert"`
			FileKey  string `yaml:"key" mapstructure:"key"`
			FileCA   string `yaml:"ca" mapstructure:"ca"`
		} `yaml:"tls" mapstructure:"tls"`
	} `yaml:"server" mapstructure:"server"`

	API struct {
		Endpoint string `yaml:"endpoint" mapstructure:"endpoint"`
	} `yaml:"api" mapstructure:"api"`

	Workdir string `yaml:"workdir" mapstructure:"workdir"`

	Manifest struct {
		Path string `yaml:"dir" mapstructure:"dir"`
	} `yaml:"manifest" mapstructure:"manifest"`

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

	fs.StringVarP(&f.AccessToken, "access-token", "", "", "Access token to API server")
	fs.StringVarP(&f.Workdir, "workdir", "", "", "Node workdir for runtime")
	fs.StringVarP(&f.Manifest.Path, "manifest-path", "", "", "Node manifest(s) path")
	fs.StringVarP(&f.API.Endpoint, "endpoint", "", "", "REST API endpoint")
}
