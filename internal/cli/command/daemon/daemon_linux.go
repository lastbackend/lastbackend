// +build linux
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
	"github.com/lastbackend/lastbackend/internal/daemon"
	"github.com/lastbackend/lastbackend/internal/daemon/config"
	"github.com/pkg/errors"
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
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			var cfg = config.Config{}

			// short-circuit on help
			help, err := cmd.Flags().GetBool("help")
			if err != nil {
				return errors.Wrapf(err, "\"help\" flag is non-bool, programmer error, please correct")
			}
			if help {
				return cmd.Help()
			}

			cfgFile, err := cmd.Flags().GetString("config")
			if err != nil {
				return errors.Wrapf(err, "\"config\" flag is non-string, programmer error, please correct")
			}

			printFlags(cmd.Flags())

			var (
				sigs = make(chan os.Signal)
			)

			if cfgFile != "" {
				if err := SetConfigFromFile(cfgFile, &cfg); err != nil {
					return errors.Wrapf(err, "can not be set server config from file")
				}
			}

			if err := SetConfigFromEnvs(&cfg); err != nil {
				return errors.Wrapf(err, "can not be set server config from envs")
			}

			if err := SetConfigFromFlags(cmd.Flags(), &cfg); err != nil {
				return errors.Wrapf(err, "can not be set server config from flags")
			}

			app, err := daemon.New(cfg)
			if err != nil {
				return errors.Wrapf(err, "can not init daemon process")
			}
			if err := app.Run(); err != nil {
				return errors.Wrapf(err, "can not run daemon process")
			}

			for {
				select {
				case <-sigs:
					app.Stop()
					return nil
				}
			}

		},
	}

	cmd.PersistentFlags().BoolP("help", "h", false, fmt.Sprintf("Help for %s", cmd.Name()))

	cmd.Flags().StringP("config", "c", "", "set config path")
	cmd.Flags().Bool("agent", false, "Only agent mode")
	cmd.Flags().Bool("no-schedule", false, "Disable schedule mode")
	cmd.Flags().StringP("access-token", "", "", "Access token to NodeClient server")
	cmd.Flags().StringP("cluster-name", "", "", "Cluster name info")
	cmd.Flags().StringP("cluster-desc", "", "", "Cluster description")
	cmd.Flags().StringP("bind-address", "", config.DefaultBindServerAddress, "Bind address for listening")
	cmd.Flags().UintP("bind-port", "", config.DefaultBindServerPort, "Bind address for listening")
	cmd.Flags().BoolP("tls-verify", "", false, "Enable check tls for NodeClient server")
	cmd.Flags().StringP("tls-cert-file", "", "", "TLS cert file path")
	cmd.Flags().StringP("tls-private-key-file", "", "", "TLS private key file path")
	cmd.Flags().StringP("tls-ca-file", "", "", "TLS certificate authority file path")
	cmd.Flags().StringP("vault-token", "", "", "Vault access token")
	cmd.Flags().StringP("vault-endpoint", "", "", "Vault access endpoint")
	cmd.Flags().StringP("domain-internal", "", config.DefaultInternalDomain, "Default external domain for cluster")
	cmd.Flags().StringP("domain-external", "", "", "Internal domain name for cluster")
	cmd.Flags().StringP("services-cidr", "", config.DefaultCIDR, "Services IP CIDR for internal IPAM service")
	cmd.Flags().StringP("root-dir", "", "", "Set root directory (Default: /var/lib/lastbackend/)")
	cmd.Flags().StringP("storage-driver", "", "", "Set storage driver (Default: overlay)")
	cmd.Flags().StringP("node-bind-address", "", config.DefaultBindServerAddress, "Set bind address for NodeClient server")
	cmd.Flags().UintP("node-bind-port", "", config.DefaultBindServerPort, "Set listening port binding for NodeClient server")
	cmd.Flags().BoolP("node-tls-verify", "", false, "Enable check tls for NodeClient server")
	cmd.Flags().StringP("node-tls-ca-file", "", "", "Set path to ca file for NodeClient server")
	cmd.Flags().StringP("node-tls-cert-file", "", "", "Set path to cert file for NodeClient server")
	cmd.Flags().StringP("node-tls-private-key-file", "", "", "Set path to key file for NodeClient server")
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
