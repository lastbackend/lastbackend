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

package daemon

import (
	"fmt"
	"os"

	"github.com/containers/storage/pkg/reexec"
	"github.com/lastbackend/lastbackend/internal/agent"
	agent_config "github.com/lastbackend/lastbackend/internal/agent/config"
	"github.com/lastbackend/lastbackend/internal/server"
	server_config "github.com/lastbackend/lastbackend/internal/server/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	if reexec.Init() {
		fmt.Println("containers storage init failed")
		os.Exit(1)
	}
}

// NewDaemonCmd entrypoint for CLI launcher
func NewCommand() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "daemon",
		Short: "Last.Backend Open-source PaaS",
		Long:  `Open-source system for automating deployment, scaling, and management of containerized applications.`,
		Run: func(cmd *cobra.Command, args []string) {

			// short-circuit on help
			help, err := cmd.Flags().GetBool("help")
			if err != nil {
				fmt.Println(`"help" flag is non-bool, programmer error, please correct`)
				return
			}
			if help {
				cmd.Help()
				return
			}

			disableMaster, err := cmd.Flags().GetBool("agent")
			if err != nil {
				fmt.Println(`"agent" flag is non-bool, programmer error, please correct`)
				return
			}
			noSchedule, err := cmd.Flags().GetBool("no-schedule")
			if err != nil {
				fmt.Println(`"no-schedule" flag is non-bool, programmer error, please correct`)
				return
			}
			cfgFile, err := cmd.Flags().GetString("config")
			if err != nil {
				fmt.Println(`"config" flag is non-string, programmer error, please correct`)
				return
			}

			PrintFlags(cmd.Flags())

			var (
				sigs = make(chan os.Signal)
			)

			if noSchedule && disableMaster {
				fmt.Println("\n#################################")
				fmt.Println("### All services was disable ###")
				fmt.Println("#################################\n")
				return
			}

			var masterApp *server.App
			var minionApp *agent.App

			if !disableMaster {

				var cfg = server_config.Config{}

				if cfgFile != "" {
					if err := SetServerConfigFromFile(cfgFile, &cfg); err != nil {
						fmt.Println("set server config from file err: ", err)
						return
					}
				}

				if err := SetServerConfigFromEnvs(&cfg); err != nil {
					fmt.Println("set server config from envs err: ", err)
					return
				}

				if err := SetServerConfigFromFlags(cmd.Flags(), &cfg); err != nil {
					fmt.Println("set server config from flags err: ", err)
					return
				}

				masterApp, err = server.New(cfg)
				if err != nil {
					fmt.Println("Create master application err: ", err)
					return
				}

				if err := masterApp.Run(); err != nil {
					fmt.Println("Run master application err: ", err)
					return
				}
			}

			if !noSchedule {

				var cfg = agent_config.Config{}

				if cfgFile != "" {
					if err := SetAgentConfigFromFile(cfgFile, &cfg); err != nil {
						fmt.Println("set agent config from file err: ", err)
						return
					}
				}

				if err := SetAgentConfigFromEnvs(&cfg); err != nil {
					fmt.Println("set agent config from envs err: ", err)
					return
				}

				if err := SetAgentConfigFromFlags(cmd.Flags(), &cfg); err != nil {
					fmt.Println("set agent config from flags err: ", err)
					return
				}

				minionApp, err = agent.New(cfg)
				if err != nil {
					fmt.Println("Create minion application err: ", err)
					return
				}

				if err := minionApp.Run(); err != nil {
					fmt.Println("Run minion application err: ", err)
					return
				}
			}

			for {
				select {
				case <-sigs:
					if !disableMaster {
						masterApp.Stop()
					}
					if !disableMaster {
						minionApp.Stop()
					}
					return
				}
			}

		},
	}

	cmd.PersistentFlags().StringP("config", "c", "", "set config path")
	cmd.PersistentFlags().Bool("agent", false, "Only agent mode")
	cmd.PersistentFlags().Bool("no-schedule", false, "Disable schedule mode")
	cmd.PersistentFlags().BoolP("help", "h", false, fmt.Sprintf("Help for %s", cmd.Name()))

	cmd.Flags().StringP("access-token", "", "", "Access token to API server")
	cmd.Flags().StringP("cluster-name", "", "", "Cluster name info")
	cmd.Flags().StringP("cluster-desc", "", "", "Cluster description")
	cmd.Flags().StringP("bind-address", "", server_config.DefaultBindServerAddress, "Bind address for listening")
	cmd.Flags().UintP("bind-port", "", server_config.DefaultBindServerPort, "Bind address for listening")
	cmd.Flags().BoolP("tls-verify", "", false, "Enable check tls for API server")
	cmd.Flags().StringP("tls-cert-file", "", "", "TLS cert file path")
	cmd.Flags().StringP("tls-private-key-file", "", "", "TLS private key file path")
	cmd.Flags().StringP("tls-ca-file", "", "", "TLS certificate authority file path")
	cmd.Flags().StringP("vault-token", "", "", "Vault access token")
	cmd.Flags().StringP("vault-endpoint", "", "", "Vault access endpoint")
	cmd.Flags().StringP("domain-internal", "", server_config.DefaultInternalDomain, "Default external domain for cluster")
	cmd.Flags().StringP("domain-external", "", "", "Internal domain name for cluster")
	cmd.Flags().StringP("services-cidr", "", agent_config.DefaultCIDR, "Services IP CIDR for internal IPAM service")
	cmd.Flags().BoolP("rootless", "", false, "Run rootless")
	cmd.Flags().StringP("workdir", "", agent_config.DefaultWorkDir, "Set directory path to hold state")
	cmd.Flags().StringP("node-bind-address", "", agent_config.DefaultBindServerAddress, "Set bind address for API server")
	cmd.Flags().UintP("node-bind-port", "", agent_config.DefaultBindServerPort, "Set listening port binding for API server")
	cmd.Flags().BoolP("node-tls-verify", "", false, "Enable check tls for API server")
	cmd.Flags().StringP("node-tls-ca-file", "", "", "Set path to ca file for API server")
	cmd.Flags().StringP("node-tls-cert-file", "", "", "Set path to cert file for API server")
	cmd.Flags().StringP("node-tls-private-key-file", "", "", "Set path to key file for API server")
	cmd.Flags().StringP("api-address", "", "", "Set endpoint for rest client")
	cmd.Flags().BoolP("api-tls-verify", "", false, "Enable check tls for rest client")
	cmd.Flags().StringP("api-tls-ca-file", "", "", "Set path to ca file for rest client")
	cmd.Flags().StringP("api-tls-cert-file", "", "", "Set path to cert file for rest client")
	cmd.Flags().StringP("api-tls-private-key-file", "", "", "Set path to key file rest client")
	cmd.Flags().StringP("manifestdir", "", "", "Set directory path to manifest")
	cmd.Flags().StringP("registry-config-path", "", "", "Registry configuration file path")
	cmd.Flags().BoolP("disable-selinux", "", false, "Disable SELinux if currently enabled")

	return cmd
}

// PrintFlags logs the flags in the flagset
func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		fmt.Println(fmt.Sprintf("FLAG: --%s=%q", flag.Name, flag.Value))
	})
}
