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

	"github.com/pkg/errors"

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
		RunE: func(cmd *cobra.Command, args []string) error {

			// short-circuit on help
			help, err := cmd.Flags().GetBool("help")
			if err != nil {
				return errors.Wrapf(err, "\"help\" flag is non-bool, programmer error, please correct")
			}
			if help {
				return cmd.Help()
			}

			disableMaster, err := cmd.Flags().GetBool("agent")
			if err != nil {
				return errors.Wrapf(err, "\"agent\" flag is non-bool, programmer error, please correct")
			}
			noSchedule, err := cmd.Flags().GetBool("no-schedule")
			if err != nil {
				return errors.Wrapf(err, "\"no-schedule\" flag is non-bool, programmer error, please correct")
			}
			cfgFile, err := cmd.Flags().GetString("config")
			if err != nil {
				return errors.Wrapf(err, "\"config\" flag is non-string, programmer error, please correct")
			}

			printFlags(cmd.Flags())

			var (
				sigs = make(chan os.Signal)
			)

			if noSchedule && disableMaster {
				fmt.Println("\n#################################")
				fmt.Println("### All services were disabled ###")
				fmt.Println("#################################\n")
				return nil
			}

			var masterApp *server.App
			var minionApp *agent.App

			if !disableMaster {

				var cfg = server_config.Config{}

				if cfgFile != "" {
					if err := SetServerConfigFromFile(cfgFile, &cfg); err != nil {
						return errors.Wrapf(err, "can not set server config from file")
					}
				}

				if err := SetServerConfigFromEnvs(&cfg); err != nil {
					return errors.Wrapf(err, "can not set server config from envs")
				}

				if err := SetServerConfigFromFlags(cmd.Flags(), &cfg); err != nil {
					return errors.Wrapf(err, "can not set server config from flags")
				}

				masterApp, err = server.New(cfg)
				if err != nil {
					return errors.Wrapf(err, "can not server initialize")
				}

				if err := masterApp.Run(); err != nil {
					return errors.Wrapf(err, "can not run server")
				}
			}

			if !noSchedule {

				var cfg = agent_config.Config{}

				if cfgFile != "" {
					if err := SetAgentConfigFromFile(cfgFile, &cfg); err != nil {
						return errors.Wrapf(err, "can not set agent config from file")
					}
				}

				if err := SetAgentConfigFromEnvs(&cfg); err != nil {
					return errors.Wrapf(err, "can not set agent config from envs")
				}

				if err := SetAgentConfigFromFlags(cmd.Flags(), &cfg); err != nil {
					return errors.Wrapf(err, "can not set agent config from flags")
				}

				minionApp, err = agent.New(cfg)
				if err != nil {
					return errors.Wrapf(err, "can not initialize agent")
				}

				if err := minionApp.Run(); err != nil {
					return errors.Wrapf(err, "can not run minion server")
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
					return nil
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
	cmd.Flags().StringP("root-dir", "", "", "Set root directory (Default: /var/lib/lastbackend/)")
	cmd.Flags().StringP("storage-driver", "", "", "Set storage driver (Default: overlay)")
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
	cmd.Flags().StringP("manifest-dir", "", "", "Set directory path to manifest")

	return cmd
}

// printFlags logs the flags in the flagset
func printFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		fmt.Println(fmt.Sprintf("FLAG: --%s=%q", flag.Name, flag.Value))
	})
}
