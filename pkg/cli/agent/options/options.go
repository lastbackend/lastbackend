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

import "github.com/spf13/pflag"

type ServerFlags struct {
	AccessToken string
	Server      struct {
		BindAddress string
		BindPort    uint
		TLS         struct {
			Verify   bool
			FileCert string
			FileKey  string
			FileCA   string
		}
	}
	API struct {
		URI      string
		BindPort uint
		TLS      struct {
			Verify   bool
			FileCert string
			FileKey  string
			FileCA   string
		}
	}
	Workdir       string
	ManifestPath  string
	BindInterface string
	Proxy         struct {
		Interface string
		Internal  string
		External  string
	}
	Network struct {
		Interface string
		Internal  string
		External  string
	}
	ContainerRuntime struct {
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
	ImageRuntime struct {
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
	fs.StringVarP(&f.BindInterface, "bind-interface", "", "", "Exporter bind network interface")
	fs.StringVarP(&f.Proxy.Interface, "network-proxy", "", "eth0", "Network proxy driver (ipvs by default)")
	fs.StringVarP(&f.Proxy.Interface, "network-proxy-iface-internal", "", "ipvs", "Network proxy internal interface binding")
	fs.StringVarP(&f.Proxy.External, "network-proxy-iface-external", "", "docker0", "Network proxy external interface binding")
	fs.StringVarP(&f.Network.Interface, "network-driver", "", "eth0", "Network driver (vxlan by default)")
	fs.StringVarP(&f.Network.Interface, "network-driver-iface-external", "", "vxlan", "Container overlay network external interface for host communication")
	fs.StringVarP(&f.Network.Interface, "network-driver-iface-internal", "", "eth0", "Container overlay network internal bridge interface for container intercommunications")
	fs.StringVarP(&f.ContainerRuntime.Driver, "container-runtime", "", "docker0", "Node container runtime")
	fs.StringVarP(&f.ContainerRuntime.Docker.Version, "container-runtime-docker-version", "", "1.38", "Set docker version for docker container runtime")
	fs.StringVarP(&f.ContainerRuntime.Docker.Host, "container-runtime-docker-host", "", "unix:///var/run/docker.sock", "Set docker host for docker container runtime")
	fs.BoolVarP(&f.ContainerRuntime.Docker.TLS.Verify, "container-runtime-docker-tls-verify", "", false, "Enable check tls for docker container runtime")
	fs.StringVarP(&f.ContainerRuntime.Docker.TLS.FileCA, "container-runtime-docker-tls-ca", "", "", "Set path to ca file for docker container runtime")
	fs.StringVarP(&f.ContainerRuntime.Docker.TLS.FileCert, "container-runtime-docker-tls-cert", "", "", "Set path to cert file for docker container runtime")
	fs.StringVarP(&f.ContainerRuntime.Docker.TLS.FileKey, "container-runtime-docker-tls-key", "", "", "Set path to key file for docker container runtime")
	fs.StringVarP(&f.ContainerRuntime.Docker.StoragePath, "container-storage-root", "", "/var/run/lastbackend", "Node container storage root")
	fs.StringVarP(&f.ImageRuntime.Driver, "container-runtime", "", "docker0", "Node container runtime")
	fs.StringVarP(&f.ImageRuntime.Docker.Version, "container-runtime-docker-version", "", "1.38", "Set docker version for docker container runtime")
	fs.StringVarP(&f.ImageRuntime.Docker.Host, "container-runtime-docker-host", "", "unix:///var/run/docker.sock", "Set docker host for docker container runtime")
	fs.BoolVarP(&f.ImageRuntime.Docker.TLS.Verify, "container-runtime-docker-tls-verify", "", false, "Enable check tls for docker container runtime")
	fs.StringVarP(&f.ImageRuntime.Docker.TLS.FileCA, "container-runtime-docker-tls-ca", "", "", "Set path to ca file for docker container runtime")
	fs.StringVarP(&f.ImageRuntime.Docker.TLS.FileCert, "container-runtime-docker-tls-cert", "", "", "Set path to cert file for docker container runtime")
	fs.StringVarP(&f.ImageRuntime.Docker.TLS.FileKey, "container-runtime-docker-tls-key", "", "", "Set path to key file for docker container runtime")
	fs.StringVarP(&f.ImageRuntime.Docker.StoragePath, "container-storage-root", "", "/var/run/lastbackend", "Node container storage root")
	fs.StringSliceVarP(&f.ContainerExtraHosts, "container-extra-hosts", "", []string{}, "Set hostname mappings for containers")
	fs.StringVarP(&f.Server.BindAddress, "bind-port", "", "0.0.0.0", "Node listening port binding")
	fs.UintVarP(&f.Server.BindPort, "tls-verify", "", 2969, "Node TLS verify options")
	fs.BoolVarP(&f.Server.TLS.Verify, "tls-verify", "", false, "Node TLS verify options")
	fs.StringVarP(&f.Server.TLS.FileCert, "tls-cert-file", "", "", "Node cert file path")
	fs.StringVarP(&f.Server.TLS.FileKey, "tls-private-key-file", "", "", "Node private key file path")
	fs.StringVarP(&f.Server.TLS.FileCA, "tls-ca-file", "", "", "Node certificate authority file path")
	fs.StringVarP(&f.API.URI, "api-uri", "", "", "REST API endpoint")
	fs.BoolVarP(&f.API.TLS.Verify, "api-tls-verify", "", false, "REST API endpoint")
	fs.StringVarP(&f.API.TLS.FileCert, "api-tls-cert-file", "", "", "REST API TLS certificate file path")
	fs.StringVarP(&f.API.TLS.FileKey, "api-tls-private-key-file", "", "false", "REST API TLS private key file path")
	fs.StringVarP(&f.API.TLS.FileCA, "api-tls-ca-file", "", "", "REST API TSL certificate authority file path")

}
