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
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

type ServerFlags struct {
	AccessToken string
	Server      struct {
		Host string
		Port uint
		TLS  struct {
			Verify   bool
			FileCert string
			FileKey  string
			FileCA   string
		}
	}
	API struct {
		Endpoint string
		TLS      struct {
			Verify   bool
			FileCert string
			FileKey  string
			FileCA   string
		}
	}
	Workdir           string
	ManifestPath      string
	ExporterInterface string
	Proxy             struct {
		Interface string
		Internal  string
		External  string
	}
	Network struct {
		Interface string
		Internal  string
		External  string
	}
	CRI struct {
		Driver string
		Docker struct {
			Version     string
			Host        string
			StoragePath string
			TLS         struct {
				Verify   bool
				FileCert string
				FileKey  string
				FileCA   string
			}
		}
	}
	IRI struct {
		Driver string
		Docker struct {
			Version     string
			Host        string
			StoragePath string
			TLS         struct {
				Verify   bool
				FileCert string
				FileKey  string
				FileCA   string
			}
		}
	}
	CSI struct {
		Path string
	}
	ContainerExtraHosts []string
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
	fs.StringVarP(&f.ManifestPath, "manifest-path", "", "", "Node manifest(s) path")
	fs.StringVarP(&f.ExporterInterface, "bind-interface", "", "eth0", "Exporter bind network interface")
	fs.StringVarP(&f.Proxy.Interface, "network-proxy", "", "ipvs", "Network proxy driver (ipvs by default)")
	fs.StringVarP(&f.Proxy.Internal, "network-proxy-iface-internal", "", "docker0", "Network proxy internal interface binding")
	fs.StringVarP(&f.Proxy.External, "network-proxy-iface-external", "", "eth0", "Network proxy external interface binding")
	fs.StringVarP(&f.Network.Interface, "network-driver", "", "vxlan", "Network driver (vxlan by default)")
	fs.StringVarP(&f.Network.Internal, "network-driver-iface-internal", "", "docker0", "Container overlay network internal bridge interface for container intercommunications")
	fs.StringVarP(&f.Network.External, "network-driver-iface-external", "", "eth0", "Container overlay network external interface for host communication")
	fs.StringVarP(&f.CRI.Driver, "container-runtime", "", "docker", "Node container runtime")
	fs.StringVarP(&f.CRI.Docker.Version, "container-runtime-docker-version", "", "1.38", "Set docker version for docker container runtime")
	fs.StringVarP(&f.CRI.Docker.Host, "container-runtime-docker-host", "", "unix:///var/run/docker.sock", "Set docker host for docker container runtime")
	fs.BoolVarP(&f.CRI.Docker.TLS.Verify, "container-runtime-docker-tls-verify", "", false, "Enable check tls for docker container runtime")
	fs.StringVarP(&f.CRI.Docker.TLS.FileCert, "container-runtime-docker-tls-cert", "", "", "Set path to cert file for docker container runtime")
	fs.StringVarP(&f.CRI.Docker.TLS.FileKey, "container-runtime-docker-tls-key", "", "", "Set path to key file for docker container runtime")
	fs.StringVarP(&f.CRI.Docker.TLS.FileCA, "container-runtime-docker-tls-ca", "", "", "Set path to ca file for docker container runtime")
	fs.StringVarP(&f.IRI.Driver, "container-image-runtime", "", "docker", "Node container image runtime")
	fs.StringVarP(&f.IRI.Docker.Version, "container-image-runtime-docker-version", "", "1.38", "Set docker version for docker container image runtime")
	fs.StringVarP(&f.IRI.Docker.Host, "container-image-runtime-docker-host", "", "unix:///var/run/docker.sock", "Set docker host for docker container image runtime")
	fs.BoolVarP(&f.IRI.Docker.TLS.Verify, "container-image-runtime-docker-tls-verify", "", false, "Enable check tls for docker container image runtime")
	fs.StringVarP(&f.IRI.Docker.TLS.FileCert, "container-image-runtime-docker-tls-cert", "", "", "Set path to cert file for docker container image runtime")
	fs.StringVarP(&f.IRI.Docker.TLS.FileKey, "container-image-runtime-docker-tls-key", "", "", "Set path to key file for docker container image runtime")
	fs.StringVarP(&f.IRI.Docker.TLS.FileCA, "container-image-runtime-docker-tls-ca", "", "", "Set path to ca file for docker container image runtime")
	fs.StringVarP(&f.CSI.Path, "container-storage-root", "", "/var/run/lastbackend", "Node container storage root")
	fs.StringSliceVarP(&f.ContainerExtraHosts, "container-extra-hosts", "", []string{}, "Set hostname mappings for containers")
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

func AddGlobalFlags(fs *pflag.FlagSet) {

	// lookup flags in global flag set and re-register the values with our flagset
	global := pflag.CommandLine
	local := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	pflagRegister(global, local, "verbose", "v")

	fs.AddFlagSet(local)
}

func pflagRegister(global, local *pflag.FlagSet, globalName string, shaortName string) {
	if f := global.Lookup(globalName); f != nil {
		f.Name = normalize(f.Name)
		f.Shorthand = shaortName
		local.AddFlag(f)
	} else {
		panic(fmt.Sprintf("failed to find flag in global flagset (pflag): %s", globalName))
	}
}

func normalize(s string) string {
	return strings.Replace(s, "_", "-", -1)
}

