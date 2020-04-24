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

type MasterFlags struct {
	Workdir            string `yaml:"workdir" mapstructure:"workdir"`
	AccessToken        string `yaml:"token" mapstructure:"token"`
	ClusterName        string `yaml:"name" mapstructure:"name"`
	ClusterDescription string `yaml:"description" mapstructure:"description"`

	Server struct {
		Host string `yaml:"host" mapstructure:"host"`
		Port uint   `yaml:"port" mapstructure:"port"`
		TLS  struct {
			FileCert string `yaml:"cert" mapstructure:"cert"`
			FileKey  string `yaml:"key" mapstructure:"key"`
			FileCA   string `yaml:"ca" mapstructure:"ca"`
		} `yaml:"tls" mapstructure:"tls"`
	} `yaml:"server" mapstructure:"server"`

	Vault struct {
		Token    string `yaml:"token" mapstructure:"token"`
		Endpoint string `yaml:"endpoint" mapstructure:"endpoint"`
	} `yaml:"vault" mapstructure:"vault"`

	Domain struct {
		Internal string `yaml:"internal" mapstructure:"internal"`
		External string `yaml:"external" mapstructure:"external"`
	} `yaml:"domain" mapstructure:"domain"`

	CIDR     string `yaml:"cidr" mapstructure:"cidr"`
	Rootless bool   `yaml:"rootless" mapstructure:"rootless"`
}

func (cfg MasterFlags) LoadViper(v *viper.Viper) *viper.Viper {
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

func NewMasterFlags() *MasterFlags {
	s := new(MasterFlags)
	return s
}

func (f *MasterFlags) AddFlags(flags *pflag.FlagSet) {

	fs := pflag.NewFlagSet("", pflag.ExitOnError)

	defer func() {
		fs.VisitAll(func(f *pflag.Flag) {
			if len(f.Deprecated) > 0 {
				f.Hidden = false
			}
		})
		flags.AddFlagSet(fs)
	}()

	fs.StringVarP(&f.AccessToken, "access-token", "", "", "Access token to API server")
	fs.StringVarP(&f.ClusterName, "cluster-name", "", "", "Cluster name info")
	fs.StringVarP(&f.ClusterDescription, "cluster-desc", "", "", "Cluster description")
	fs.StringVarP(&f.Server.Host, "bind-address", "", "0.0.0.0", "Bind address for listening")
	fs.UintVarP(&f.Server.Port, "bind-port", "", 2967, "Bind address for listening")
	fs.StringVarP(&f.Server.TLS.FileCert, "tls-cert-file", "", "", "TLS cert file path")
	fs.StringVarP(&f.Server.TLS.FileKey, "tls-private-key-file", "", "", "TLS private key file path")
	fs.StringVarP(&f.Server.TLS.FileCA, "tls-ca-file", "", "", "TLS certificate authority file path")
	fs.StringVarP(&f.Vault.Token, "vault-token", "", "", "Vault access token")
	fs.StringVarP(&f.Vault.Endpoint, "vault-endpoint", "", "", "Vault access endpoint")
	fs.StringVarP(&f.Domain.Internal, "domain-internal", "", "lb.local", "Default external domain for cluster")
	fs.StringVarP(&f.Domain.External, "domain-external", "", "", "Internal domain name for cluster")
	fs.StringVarP(&f.CIDR, "services-cidr", "", "172.0.0.0/24", "Services IP CIDR for internal IPAM service")
	fs.StringVarP(&f.Workdir, "workdir", "", "/${HOME}/.lastbackend/", "Set directory path to hold state")
	fs.BoolVarP(&f.Rootless, "rootless", "", false, "Run rootless")

}
