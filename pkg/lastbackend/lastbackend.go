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

package lastbackend

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/agent"
	"os"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/api"
	"github.com/lastbackend/lastbackend/pkg/controller"
	"github.com/lastbackend/lastbackend/pkg/discovery"
	"github.com/lastbackend/lastbackend/pkg/exporter"
	"github.com/lastbackend/lastbackend/pkg/ingress"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultEnvPrefix = "LB"
const defaultConfigType = "yaml"
const defaultConfigName = "config"

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "lb",
	Short: "Last.Backend Open-source API",
	Long:  `Open-source system for automating deployment, scaling, and management of containerized applications.`,
	Run: func(cmd *cobra.Command, args []string) {

		var (
			done    = make(chan bool, 1)
			apps    = make(chan bool)
			wait    = 0
			daemons = []func(*viper.Viper) bool{
				api.Daemon,
				controller.Daemon,
			}
		)

		for _, app := range daemons {
			go func() {
				wait++
				apps <- app(viper.GetViper())
			}()
		}

		go func() {
			for {
				select {
				case <-apps:
					wait--
					if wait == 0 {
						done <- true
						return
					}
				}
			}
		}()

		<-done

	},
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Run node agent",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Error(err.Error())
			return
		}

		agent.Daemon(viper.GetViper())
	},
}

var ingressCmd = &cobra.Command{
	Use:   "ingress",
	Short: "Ingress instance",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Error(err.Error())
			return
		}

		ingress.Daemon(viper.GetViper())
	},
}

var discoveryCmd = &cobra.Command{
	Use:   "discovery",
	Short: "Discovery instance",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Error(err.Error())
			return
		}

		discovery.Daemon(viper.GetViper())
	},
}

var exporterCmd = &cobra.Command{
	Use:   "exporter",
	Short: "Exporter instance",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Error(err.Error())
			return
		}

		exporter.Daemon(viper.GetViper())
	},
}

func Execute() {

	cobra.OnInitialize(initConfig)

	// ===================================================================================================================
	// Root flags
	// ===================================================================================================================

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Path for the configuration file")
	rootCmd.PersistentFlags().IntP("verbose", "v", 0, "Set log level from 0 to 7")

	viper.BindPFlag("service.cidr", rootCmd.PersistentFlags().Lookup("services-cidr"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.Flags().StringP("access-token", "", "", "Access token to API server")
	rootCmd.Flags().StringP("cluster-name", "", "", "Cluster name info")
	rootCmd.Flags().StringP("cluster-description", "", "", "Cluster description")
	rootCmd.Flags().StringP("bind-address", "", "0.0.0.0", "Bind address for listening")
	rootCmd.Flags().UintP("bind-port", "", 2967, "Bind address for listening")
	rootCmd.Flags().StringP("tls-cert-file", "", "", "TLS cert file path")
	rootCmd.Flags().StringP("tls-private-key-file", "", "", "TLS private key file path")
	rootCmd.Flags().StringP("tls-ca-file", "", "", "TLS certificate authority file path")
	rootCmd.Flags().StringP("vault-token", "", "", "Vault access token")
	rootCmd.Flags().StringP("vault-endpoint", "", "", "Vault access endpoint")
	rootCmd.Flags().StringP("domain-internal", "", "lb.local", "Default external domain for cluster")
	rootCmd.Flags().StringP("domain-external", "", "", "Internal domain name for cluster")
	rootCmd.Flags().StringP("storage", "", "etcd", "Set storage driver (Allow: etcd, mock)")
	rootCmd.Flags().StringP("etcd-cert-file", "", "", "ETCD database cert file path")
	rootCmd.Flags().StringP("etcd-private-key-file", "", "", "ETCD database private key file path")
	rootCmd.Flags().StringP("etcd-ca-file", "", "", "ETCD database certificate authority file")
	rootCmd.Flags().StringSliceP("etcd-endpoints", "", []string{"127.0.0.1:2379"}, "ETCD database endpoints list")
	rootCmd.Flags().StringP("services-cidr", "", "172.0.0.0/24", "Services IP CIDR for internal IPAM service")

	rootCmd.MarkFlagRequired("vault-token")
	rootCmd.MarkFlagRequired("vault-endpoint")

	viper.BindPFlag("token", rootCmd.Flags().Lookup("access-token"))
	viper.BindPFlag("name", rootCmd.Flags().Lookup("cluster-name"))
	viper.BindPFlag("description", rootCmd.Flags().Lookup("cluster-description"))
	viper.BindPFlag("server.host", rootCmd.Flags().Lookup("bind-address"))
	viper.BindPFlag("server.port", rootCmd.Flags().Lookup("bind-port"))
	viper.BindPFlag("server.tls.cert", rootCmd.Flags().Lookup("tls-cert-file"))
	viper.BindPFlag("server.tls.key", rootCmd.Flags().Lookup("tls-private-key-file"))
	viper.BindPFlag("server.tls.ca", rootCmd.Flags().Lookup("tls-ca-file"))
	viper.BindPFlag("vault.token", rootCmd.Flags().Lookup("vault-token"))
	viper.BindPFlag("vault.endpoint", rootCmd.Flags().Lookup("vault-endpoint"))
	viper.BindPFlag("domain.internal", rootCmd.Flags().Lookup("domain-internal"))
	viper.BindPFlag("domain.external", rootCmd.Flags().Lookup("domain-external"))
	viper.BindPFlag("storage.driver", rootCmd.Flags().Lookup("storage"))
	viper.BindPFlag("storage.etcd.tls.cert", rootCmd.Flags().Lookup("etcd-cert-file"))
	viper.BindPFlag("storage.etcd.tls.key", rootCmd.Flags().Lookup("etcd-private-key-file"))
	viper.BindPFlag("storage.etcd.tls.ca", rootCmd.Flags().Lookup("etcd-ca-file"))
	viper.BindPFlag("storage.etcd.endpoints", rootCmd.Flags().Lookup("etcd-endpoints"))
	viper.BindPFlag("storage.etcd.prefix", rootCmd.Flags().Lookup("etcd-prefix"))

	// ===================================================================================================================
	// Node agent flags
	// ===================================================================================================================

	agentCmd.Flags().StringP("access-token", "", "", "Access token to API server")
	agentCmd.Flags().StringP("workdir", "", "", "Node workdir for runtime")
	agentCmd.Flags().StringP("manifest-path", "", "", "Node manifest(s) path")
	agentCmd.Flags().StringP("bind-interface", "", "", "Exporter bind network interface")
	agentCmd.Flags().StringP("network-proxy", "", "eth0", "Network proxy driver (ipvs by default)")
	agentCmd.Flags().StringP("network-proxy-iface-internal", "", "ipvs", "Network proxy internal interface binding")
	agentCmd.Flags().StringP("network-proxy-iface-external", "", "docker0", "Network proxy external interface binding")
	agentCmd.Flags().StringP("network-driver", "", "eth0", "Network driver (vxlan by default)")
	agentCmd.Flags().StringP("network-driver-iface-external", "", "vxlan", "Container overlay network external interface for host communication")
	agentCmd.Flags().StringP("network-driver-iface-internal", "", "eth0", "Container overlay network internal bridge interface for container intercommunications")
	agentCmd.Flags().StringP("container-runtime", "", "docker0", "Node container runtime")
	agentCmd.Flags().StringP("container-runtime-docker-version", "", "docker", "Set docker version for docker container runtime")
	agentCmd.Flags().StringP("container-runtime-docker-host", "", "1.38", "Set docker host for docker container runtime")
	agentCmd.Flags().StringP("container-runtime-docker-tls-verify", "", "unix:///var/run/docker.sock", "Enable check tls for docker container runtime")
	agentCmd.Flags().BoolP("container-runtime-docker-tls-ca", "", false, "Set path to ca file for docker container runtime")
	agentCmd.Flags().StringP("container-runtime-docker-tls-cert", "", "", "Set path to cert file for docker container runtime")
	agentCmd.Flags().StringP("container-runtime-docker-tls-key", "", "", "Set path to key file for docker container runtime")
	agentCmd.Flags().StringP("container-storage-root", "", "", "Node container storage root")
	agentCmd.Flags().StringP("container-image-runtime", "", "/var/run/lastbackend", "Node container images runtime")
	agentCmd.Flags().StringP("container-image-runtime-docker-version", "", "docker", "Set docker version for docker container image runtime")
	agentCmd.Flags().StringP("container-image-runtime-docker-host", "", "1.38", "Set docker host for docker container image runtime")
	agentCmd.Flags().StringP("container-image-runtime-docker-tls-verify", "", "unix:///var/run/docker.sock", "Enable check tls for docker container image runtime")
	agentCmd.Flags().BoolP("container-image-runtime-docker-tls-ca", "", false, "Set path to ca file for docker container image runtime")
	agentCmd.Flags().StringP("container-image-runtime-docker-tls-cert", "", "", "Set path to cert file for docker container image runtime")
	agentCmd.Flags().StringP("container-image-runtime-docker-tls-key", "", "", "Set path to key file for docker container image runtime")
	agentCmd.Flags().StringP("container-extra-hosts", "", "", "Set hostname mappings for containers")
	agentCmd.Flags().StringSliceP("bind-address", "", []string{}, "Node bind address")
	agentCmd.Flags().StringP("bind-port", "", "0.0.0.0", "Node listening port binding")
	agentCmd.Flags().UintP("tls-verify", "", 2969, "Node TLS verify options")
	agentCmd.Flags().BoolP("tls-cert-file", "", false, "Node cert file path")
	agentCmd.Flags().StringP("tls-private-key-file", "", "", "Node private key file path")
	agentCmd.Flags().StringP("tls-ca-file", "", "", "Node certificate authority file path")
	agentCmd.Flags().StringP("api-uri", "", "", "REST API endpoint")
	agentCmd.Flags().StringP("api-tls-verify", "", "", "REST API endpoint")
	agentCmd.Flags().StringP("api-tls-cert-file", "", "", "REST API TLS certificate file path")
	agentCmd.Flags().BoolP("api-tls-private-key-file", "", false, "REST API TLS private key file path")
	agentCmd.Flags().StringP("api-tls-ca-file", "", "", "REST API TSL certificate authority file path")

	viper.BindPFlag("token", agentCmd.Flags().Lookup("access-token"))
	viper.BindPFlag("workdir", agentCmd.Flags().Lookup("workdir"))
	viper.BindPFlag("manifest.dir", agentCmd.Flags().Lookup("manifest-path"))
	viper.BindPFlag("network.interface", agentCmd.Flags().Lookup("bind-interface"))
	viper.BindPFlag("network.cpi.type", agentCmd.Flags().Lookup("network-proxy"))
	viper.BindPFlag("network.cpi.interface.internal", agentCmd.Flags().Lookup("network-proxy-iface-internal"))
	viper.BindPFlag("network.cpi.interface.external", agentCmd.Flags().Lookup("network-proxy-iface-external"))
	viper.BindPFlag("network.cni.type", agentCmd.Flags().Lookup("network-driver"))
	viper.BindPFlag("network.cni.interface.external", agentCmd.Flags().Lookup("network-driver-iface-external"))
	viper.BindPFlag("network.cni.interface.internal", agentCmd.Flags().Lookup("network-driver-iface-internal"))
	viper.BindPFlag("container.cri.type", agentCmd.Flags().Lookup("container-runtime"))
	viper.BindPFlag("container.cri.docker.version", agentCmd.Flags().Lookup("container-runtime-docker-version"))
	viper.BindPFlag("container.cri.docker.host", agentCmd.Flags().Lookup("container-runtime-docker-host"))
	viper.BindPFlag("container.cri.docker.tls.verify", agentCmd.Flags().Lookup("container-runtime-docker-tls-verify"))
	viper.BindPFlag("container.cri.docker.tls.ca_file", agentCmd.Flags().Lookup("container-runtime-docker-tls-ca"))
	viper.BindPFlag("container.cri.docker.tls.cert_file", agentCmd.Flags().Lookup("container-runtime-docker-tls-cert"))
	viper.BindPFlag("container.cri.docker.tls.key_file", agentCmd.Flags().Lookup("container-runtime-docker-tls-key"))
	viper.BindPFlag("container.csi.dir.root", agentCmd.Flags().Lookup("container-storage-root"))
	viper.BindPFlag("container.iri.type", agentCmd.Flags().Lookup("container-image-runtime"))
	viper.BindPFlag("container.iri.docker.version", agentCmd.Flags().Lookup("container-image-runtime-docker-version"))
	viper.BindPFlag("container.iri.docker.host", agentCmd.Flags().Lookup("container-image-runtime-docker-host"))
	viper.BindPFlag("container.iri.docker.tls.verify", agentCmd.Flags().Lookup("container-image-runtime-docker-tls-verify"))
	viper.BindPFlag("container.iri.docker.tls.ca_file", agentCmd.Flags().Lookup("container-image-runtime-docker-tls-ca"))
	viper.BindPFlag("container.iri.docker.tls.cert_file", agentCmd.Flags().Lookup("container-image-runtime-docker-tls-cert"))
	viper.BindPFlag("container.iri.docker.tls.key_file", agentCmd.Flags().Lookup("container-image-runtime-docker-tls-key"))
	viper.BindPFlag("container.extra_hosts", agentCmd.Flags().Lookup("container-extra-hosts"))
	viper.BindPFlag("server.host", agentCmd.Flags().Lookup("bind-address"))
	viper.BindPFlag("server.port", agentCmd.Flags().Lookup("bind-port"))
	viper.BindPFlag("server.tls.verify", agentCmd.Flags().Lookup("tls-verify"))
	viper.BindPFlag("server.tls.cert", agentCmd.Flags().Lookup("tls-cert-file"))
	viper.BindPFlag("server.tls.key", agentCmd.Flags().Lookup("tls-private-key-file"))
	viper.BindPFlag("server.tls.ca", agentCmd.Flags().Lookup("tls-ca-file"))
	viper.BindPFlag("api.uri", agentCmd.Flags().Lookup("api-uri"))
	viper.BindPFlag("api.tls.verify", agentCmd.Flags().Lookup("api-tls-verify"))
	viper.BindPFlag("api.tls.cert", agentCmd.Flags().Lookup("api-tls-cert-file"))
	viper.BindPFlag("api.tls.key", agentCmd.Flags().Lookup("api-tls-private-key-file"))
	viper.BindPFlag("api.tls.ca", agentCmd.Flags().Lookup("api-tls-ca-file"))

	rootCmd.AddCommand(agentCmd)

	// ===================================================================================================================
	// Ingress flags
	// ===================================================================================================================

	ingressCmd.Flags().StringP("access-token", "", "", "Access token to API server")
	ingressCmd.Flags().StringP("haproxy-config-path", "", "/var/run/lastbackend/ingress/haproxy", "HAProxy configuration path setup")
	ingressCmd.Flags().StringP("haproxy-pid", "", "/var/run/lastbackend/ingress/haproxy/haproxy.pid", "HAProxy pid file path")
	ingressCmd.Flags().StringP("haproxy-exec", "", "/usr/sbin/haproxy", "HAProxy entrypoint path")
	ingressCmd.Flags().UintP("haproxy-stat-port", "", 1936, "HAProxy statistic port definition. If not provided - statistic will be disabled")
	ingressCmd.Flags().StringP("haproxy-stat-username", "", "", "HAProxy statistic access username")
	ingressCmd.Flags().StringP("haproxy-stat-password", "", "", "HAProxy statistic access password")
	ingressCmd.Flags().StringP("bind-interface", "", "eth0", "Exporter bind network interface")
	ingressCmd.Flags().StringP("network-proxy", "", "ipvs", "Container proxy interface driver")
	ingressCmd.Flags().StringP("network-proxy-iface-internal", "", "docker0", "Network external interface binding")
	ingressCmd.Flags().StringP("network-proxy-iface-external", "", "eth0", "Network container bridge binding")
	ingressCmd.Flags().StringP("network-driver", "", "vxlan", "Container overlay network driver")
	ingressCmd.Flags().StringP("network-driver-iface-external", "", "eth0", "")
	ingressCmd.Flags().StringP("network-driver-iface-bridge", "", "docker0", "")
	ingressCmd.Flags().StringSliceP("network-resolvers", "", []string{"8.8.8.8", "8.8.4.4"}, "Additional resolvers IPS for Ingress")
	ingressCmd.Flags().StringP("api-uri", "", "", "REST API endpoint")
	ingressCmd.Flags().StringP("api-cert-file", "", "", "REST API TLS certificate file path")
	ingressCmd.Flags().StringP("api-private-key-file", "", "", "REST API TLS private key file path")
	ingressCmd.Flags().StringP("api-ca-file", "", "", "REST API TSL certificate authority file path")

	viper.BindPFlag("token", ingressCmd.Flags().Lookup("access-token"))
	viper.BindPFlag("haproxy.config", ingressCmd.Flags().Lookup("haproxy-config-path"))
	viper.BindPFlag("haproxy.pid", ingressCmd.Flags().Lookup("haproxy-pid"))
	viper.BindPFlag("haproxy.exec", ingressCmd.Flags().Lookup("haproxy-exec"))
	viper.BindPFlag("haproxy.stat.port", ingressCmd.Flags().Lookup("haproxy-stat-port"))
	viper.BindPFlag("haproxy.stat.username", ingressCmd.Flags().Lookup("haproxy-stat-username"))
	viper.BindPFlag("haproxy.stat.password", ingressCmd.Flags().Lookup("haproxy-stat-password"))
	viper.BindPFlag("network.interface", ingressCmd.Flags().Lookup("bind-interface"))
	viper.BindPFlag("network.cpi.type", ingressCmd.Flags().Lookup("network-proxy"))
	viper.BindPFlag("network.cpi.interface.internal", ingressCmd.Flags().Lookup("network-proxy-iface-internal"))
	viper.BindPFlag("network.cpi.interface.external", ingressCmd.Flags().Lookup("network-proxy-iface-external"))
	viper.BindPFlag("network.cni.type", ingressCmd.Flags().Lookup("network-driver"))
	viper.BindPFlag("network.cni.interface.external", ingressCmd.Flags().Lookup("network-driver-iface-external"))
	viper.BindPFlag("network.cni.interface.internal", ingressCmd.Flags().Lookup("network-driver-iface-bridge"))
	viper.BindPFlag("resolver.servers", ingressCmd.Flags().Lookup("network-resolvers"))
	viper.BindPFlag("api.uri", ingressCmd.Flags().Lookup("api-uri"))
	viper.BindPFlag("api.tls.cert", ingressCmd.Flags().Lookup("api-cert-file"))
	viper.BindPFlag("api.tls.key", ingressCmd.Flags().Lookup("api-private-key-file"))
	viper.BindPFlag("api.tls.ca", ingressCmd.Flags().Lookup("api-ca-file"))

	rootCmd.AddCommand(ingressCmd)

	// ===================================================================================================================
	// Discovery flags
	// ===================================================================================================================

	discoveryCmd.Flags().StringP("access-token", "", "", "Access token to API server")
	discoveryCmd.Flags().StringP("storage", "", "etcd", "Set storage driver (Allow: etcd, mock)")
	discoveryCmd.Flags().StringSliceP("etcd-endpoints", "", []string{"127.0.0.1:2379"}, "ETCD database endpoints list")
	discoveryCmd.Flags().StringP("etcd-prefix", "", "lastbackend", "ETCD database storage prefix")
	discoveryCmd.Flags().StringP("etcd-cert-file", "", "", "ETCD database cert file path")
	discoveryCmd.Flags().StringP("etcd-private-key-file", "", "", "ETCD database private key file path")
	discoveryCmd.Flags().StringP("etcd-ca-file", "", "", "ETCD database certificate authority file")
	discoveryCmd.Flags().StringP("bind-address", "", "0.0.0.0", "DNS server bind address")
	discoveryCmd.Flags().UintP("bind-port", "", 53, "DNS port listening")
	discoveryCmd.Flags().StringP("dns-ttl", "", "24h", "DNS cache ttl")
	discoveryCmd.Flags().StringP("api-uri", "", "", "REST API endpoint")
	discoveryCmd.Flags().StringP("api-cert-file", "", "", "REST API TLS certificate file path")
	discoveryCmd.Flags().StringP("api-private-key-file", "", "", "REST API TLS private key file path")
	discoveryCmd.Flags().StringP("api-ca-file", "", "", "REST API TSL certificate authority file path")

	viper.BindPFlag("token", discoveryCmd.Flags().Lookup("access-token"))
	viper.BindPFlag("storage.driver", discoveryCmd.Flags().Lookup("storage"))
	viper.BindPFlag("storage.etcd.endpoints", discoveryCmd.Flags().Lookup("etcd-endpoints"))
	viper.BindPFlag("storage.etcd.prefix", discoveryCmd.Flags().Lookup("etcd-prefix"))
	viper.BindPFlag("storage.etcd.tls.cert", discoveryCmd.Flags().Lookup("etcd-cert-file"))
	viper.BindPFlag("storage.etcd.tls.key", discoveryCmd.Flags().Lookup("etcd-private-key-file"))
	viper.BindPFlag("storage.etcd.tls.ca", discoveryCmd.Flags().Lookup("etcd-ca-file"))
	viper.BindPFlag("dns.host", discoveryCmd.Flags().Lookup("bind-address"))
	viper.BindPFlag("dns.port", discoveryCmd.Flags().Lookup("bind-port"))
	viper.BindPFlag("dns.ttl", discoveryCmd.Flags().Lookup("dns-ttl"))
	viper.BindPFlag("api.uri", discoveryCmd.Flags().Lookup("api-uri"))
	viper.BindPFlag("api.tls.cert", discoveryCmd.Flags().Lookup("api-cert-file"))
	viper.BindPFlag("api.tls.key", discoveryCmd.Flags().Lookup("api-private-key-file"))
	viper.BindPFlag("api.tls.ca", discoveryCmd.Flags().Lookup("api-ca-file"))

	rootCmd.AddCommand(discoveryCmd)

	// ===================================================================================================================
	// Exporter flags
	// ===================================================================================================================

	exporterCmd.Flags().StringP("access-token", "", "", "Access token to API server")
	exporterCmd.Flags().StringP("bind-listener-address", "", "0.0.0.0", "Exporter bind address")
	exporterCmd.Flags().UintP("bind-listener-port", "", 2963, "Exporter bind port")
	exporterCmd.Flags().StringP("bind-rest-address", "", "0.0.0.0", "Exporter REST listener address")
	exporterCmd.Flags().UintP("bind-rest-port", "", 2964, "Exporter REST listener port")
	exporterCmd.Flags().StringP("tls-cert-file", "", "", "Exporter REST TLS cert file path")
	exporterCmd.Flags().StringP("tls-private-key-file", "", "", "Exportter REST TLS private key path")
	exporterCmd.Flags().StringP("tls-ca-file", "", "", "Exporter REST TLS certificate authority file path")
	exporterCmd.Flags().StringP("api-uri", "", "", "REST API endpoint")
	exporterCmd.Flags().StringP("api-cert-file", "", "", "REST API TLS certificate file path")
	exporterCmd.Flags().StringP("api-private-key-file", "", "", "REST API TLS private key file path")
	exporterCmd.Flags().StringP("api-ca-file", "", "", "REST API TSL certificate authority file path")
	exporterCmd.Flags().StringP("bind-interface", "", "eth0", "Exporter bind network interface")
	exporterCmd.Flags().StringP("log-workdir", "", "/var/run/lastbackend", "Set directory on host for logs storage")

	viper.BindPFlag("token", exporterCmd.Flags().Lookup("access-token"))
	viper.BindPFlag("logger.host", exporterCmd.Flags().Lookup("bind-listener-address"))
	viper.BindPFlag("logger.port", exporterCmd.Flags().Lookup("bind-listener-port"))
	viper.BindPFlag("server.host", exporterCmd.Flags().Lookup("bind-rest-address"))
	viper.BindPFlag("server.port", exporterCmd.Flags().Lookup("bind-rest-port"))
	viper.BindPFlag("server.tls.cert", exporterCmd.Flags().Lookup("tls-cert-file"))
	viper.BindPFlag("server.tls.key", exporterCmd.Flags().Lookup("tls-private-key-file"))
	viper.BindPFlag("server.tls.ca", exporterCmd.Flags().Lookup("tls-ca-file"))
	viper.BindPFlag("api.uri", exporterCmd.Flags().Lookup("api-uri"))
	viper.BindPFlag("api.tls.cert", exporterCmd.Flags().Lookup("api-cert-file"))
	viper.BindPFlag("api.tls.key", exporterCmd.Flags().Lookup("api-private-key-file"))
	viper.BindPFlag("api.tls.ca", exporterCmd.Flags().Lookup("api-ca-file"))
	viper.BindPFlag("network.interface", exporterCmd.Flags().Lookup("bind-interface"))
	viper.BindPFlag("logger.workdir", exporterCmd.Flags().Lookup("log-workdir"))

	rootCmd.AddCommand(exporterCmd)

	// ===================================================================================================================
	// Execute commands
	// ===================================================================================================================

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix(defaultEnvPrefix)
	viper.SetConfigType(defaultConfigType)
	viper.SetConfigFile(viper.GetString(defaultConfigName))

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Can't read config:", err)
			os.Exit(1)
		}
	}
}
