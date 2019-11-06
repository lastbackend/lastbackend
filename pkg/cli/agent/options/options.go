//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

type ServerFlags struct {
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
		Endpoint string `yaml:"uri" mapstructure:"uri"`
		TLS      struct {
			Verify   bool   `yaml:"verify" mapstructure:"verify"`
			FileCert string `yaml:"cert" mapstructure:"cert"`
			FileKey  string `yaml:"key" mapstructure:"key"`
			FileCA   string `yaml:"ca" mapstructure:"ca"`
		} `yaml:"tls" mapstructure:"tls"`
	} `yaml:"api" mapstructure:"api"`

	Workdir string `yaml:"workdir" mapstructure:"workdir"`

	Manifest struct {
		Path string `yaml:"dir" mapstructure:"dir"`
	} `yaml:"manifest" mapstructure:"manifest"`

	Network struct {
		Type string `yaml:"interface" mapstructure:"interface"`
		CPI  struct {
			Type      string `yaml:"type" mapstructure:"type"`
			Interface struct {
				Internal string `yaml:"internal" mapstructure:"internal"`
				External string `yaml:"external" mapstructure:"external"`
			} `yaml:"interface" mapstructure:"interface"`
		} `yaml:"cpi" mapstructure:"cpi"`

		CNI struct {
			Type      string `yaml:"type" mapstructure:"type"`
			Interface struct {
				Internal string `yaml:"internal" mapstructure:"internal"`
				External string `yaml:"external" mapstructure:"external"`
			} `yaml:"interface" mapstructure:"interface"`
		} `yaml:"cni" mapstructure:"cni"`
	} `yaml:"network" mapstructure:"network"`

	Container struct {
		CRI struct {
			Driver string `yaml:"type" mapstructure:"type"`
			Docker struct {
				Version string `yaml:"version" mapstructure:"version"`
				Host    string `yaml:"host" mapstructure:"host"`
				TLS     struct {
					Verify   bool   `yaml:"verify" mapstructure:"verify"`
					FileCert string `yaml:"cert_file" mapstructure:"cert_file"`
					FileKey  string `yaml:"key_file" mapstructure:"key_file"`
					FileCA   string `yaml:"ca_file" mapstructure:"ca_file"`
				} `yaml:"tls" mapstructure:"tls"`
			} `yaml:"docker" mapstructure:"docker"`
		} `yaml:"cri" mapstructure:"cri"`

		IRI struct {
			Driver string `yaml:"type" mapstructure:"type"`
			Docker struct {
				Version string `yaml:"version" mapstructure:"version"`
				Host    string `yaml:"host" mapstructure:"host"`
				TLS     struct {
					Verify   bool   `yaml:"verify" mapstructure:"verify"`
					FileCert string `yaml:"cert_file" mapstructure:"cert_file"`
					FileKey  string `yaml:"key_file" mapstructure:"key_file"`
					FileCA   string `yaml:"ca_file" mapstructure:"ca_file"`
				} `yaml:"tls" mapstructure:"tls"`
			} `yaml:"docker" mapstructure:"docker"`
		} `yaml:"iri" mapstructure:"iri"`

		CSI struct {
			Dir struct {
				RootPath string `yaml:"root" mapstructure:"root"`
			} `yaml:"dir" mapstructure:"dir"`
		} `yaml:"csi" mapstructure:"csi"`
		ContainerExtraHosts []string `yaml:"extra_hosts" mapstructure:"extra_hosts"`
	} `yaml:"container" mapstructure:"container"`
}

func (cfg *ServerFlags) LoadViper(v *viper.Viper) *viper.Viper {
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

func NewServerFlags() *ServerFlags {
	s := new(ServerFlags)
	return s
}

func (f *ServerFlags) AddFlags(mainfs *pflag.FlagSet) {

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
	fs.StringVarP(&f.Network.Type, "bind-interface", "", "eth0", "Exporter bind network interface")
	fs.StringVarP(&f.Network.CPI.Type, "network-proxy", "", "ipvs", "Network proxy driver (ipvs by default)")
	fs.StringVarP(&f.Network.CPI.Interface.Internal, "network-proxy-iface-internal", "", "docker0", "Network proxy internal interface binding")
	fs.StringVarP(&f.Network.CPI.Interface.External, "network-proxy-iface-external", "", "eth0", "Network proxy external interface binding")
	fs.StringVarP(&f.Network.CNI.Type, "network-driver", "", "vxlan", "Network driver (vxlan by default)")
	fs.StringVarP(&f.Network.CNI.Interface.Internal, "network-driver-iface-internal", "", "docker0", "Container overlay network internal bridge interface for container intercommunications")
	fs.StringVarP(&f.Network.CNI.Interface.External, "network-driver-iface-external", "", "eth0", "Container overlay network external interface for host communication")
	fs.StringVarP(&f.Container.CRI.Driver, "container-runtime", "", "docker", "Node container runtime")
	fs.StringVarP(&f.Container.CRI.Docker.Version, "container-runtime-docker-version", "", "1.38", "Set docker version for docker container runtime")
	fs.StringVarP(&f.Container.CRI.Docker.Host, "container-runtime-docker-host", "", "unix:///var/run/docker.sock", "Set docker host for docker container runtime")
	fs.BoolVarP(&f.Container.CRI.Docker.TLS.Verify, "container-runtime-docker-tls-verify", "", false, "Enable check tls for docker container runtime")
	fs.StringVarP(&f.Container.CRI.Docker.TLS.FileCert, "container-runtime-docker-tls-cert", "", "", "Set path to cert file for docker container runtime")
	fs.StringVarP(&f.Container.CRI.Docker.TLS.FileKey, "container-runtime-docker-tls-key", "", "", "Set path to key file for docker container runtime")
	fs.StringVarP(&f.Container.CRI.Docker.TLS.FileCA, "container-runtime-docker-tls-ca", "", "", "Set path to ca file for docker container runtime")
	fs.StringVarP(&f.Container.IRI.Driver, "container-image-runtime", "", "docker", "Node container image runtime")
	fs.StringVarP(&f.Container.IRI.Docker.Version, "container-image-runtime-docker-version", "", "1.38", "Set docker version for docker container image runtime")
	fs.StringVarP(&f.Container.IRI.Docker.Version, "container-image-runtime-docker-host", "", "unix:///var/run/docker.sock", "Set docker host for docker container image runtime")
	fs.BoolVarP(&f.Container.IRI.Docker.TLS.Verify, "container-image-runtime-docker-tls-verify", "", false, "Enable check tls for docker container image runtime")
	fs.StringVarP(&f.Container.IRI.Docker.TLS.FileCert, "container-image-runtime-docker-tls-cert", "", "", "Set path to cert file for docker container image runtime")
	fs.StringVarP(&f.Container.IRI.Docker.TLS.FileKey, "container-image-runtime-docker-tls-key", "", "", "Set path to key file for docker container image runtime")
	fs.StringVarP(&f.Container.IRI.Docker.TLS.FileCA, "container-image-runtime-docker-tls-ca", "", "", "Set path to ca file for docker container image runtime")
	fs.StringVarP(&f.Container.CSI.Dir.RootPath, "container-storage-root", "", "/var/run/lastbackend", "Node container storage root")
	fs.StringSliceVarP(&f.Container.ContainerExtraHosts, "container-extra-hosts", "", []string{}, "Set hostname mappings for containers")
	fs.StringVarP(&f.Server.Host, "bind-address", "", "0.0.0.0", "Node bind address")
	fs.UintVarP(&f.Server.Port, "bind-port", "", 2969, "Node listening port binding")
	fs.BoolVarP(&f.Server.TLS.Verify, "tls-verify", "", false, "Node TLS verify options")
	fs.StringVarP(&f.Server.TLS.FileCert, "tls-cert-file", "", "", "Node cert file path")
	fs.StringVarP(&f.Server.TLS.FileKey, "tls-private-key-file", "", "", "Node private key file path")
	fs.StringVarP(&f.Server.TLS.FileCA, "tls-ca-file", "", "", "Node certificate authority file path")
	fs.StringVarP(&f.API.Endpoint, "api-uri", "", "", "REST API endpoint")
	fs.BoolVarP(&f.API.TLS.Verify, "api-tls-verify", "", false, "REST API TLS verify options")
	fs.StringVarP(&f.API.TLS.FileCert, "api-tls-cert-file", "", "", "REST API TLS certificate file path")
	fs.StringVarP(&f.API.TLS.FileKey, "api-tls-private-key-file", "", "false", "REST API TLS private key file path")
	fs.StringVarP(&f.API.TLS.FileCA, "api-tls-ca-file", "", "", "REST API TSL certificate authority file path")
}
