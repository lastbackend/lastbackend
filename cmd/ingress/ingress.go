//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/ingress"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config string //
	token  string //
	daemon bool   //
	debug  int    //

	// HoarderCmd ...
	CLI = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  ``,

		// parse the config if one is provided, or use the defaults. Set the backend
		// driver to be used
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

			// if --config is passed, attempt to parse the config file
			if config != "" {

				// get the filepath
				abs, err := filepath.Abs(config)
				if err != nil {
					fmt.Printf("Error reading filepath: %s \n", err)
				}

				// get the config name
				base := filepath.Base(abs)

				// get the path
				path := filepath.Dir(abs)

				//
				viper.SetConfigName(strings.Split(base, ".")[0])
				viper.AddConfigPath(path)

				// Find and read the config file; Handle errors reading the config file
				if err := viper.ReadInConfig(); err != nil {
					fmt.Printf("Failed to read config file: %s\n", err)
					os.Exit(1)
				}
			}
		},

		// either run hoarder as a server, or run it as a CLI depending on what flags
		// are provided
		Run: func(cmd *cobra.Command, args []string) {

			// if --server is passed start the hoarder server
			if daemon {
				// do server stuff...
			}

			ingress.Daemon()
		},
	}
)

func init() {

	// set config defaults
	viper.SetDefault("garbage-collect", false)

	// local flags;
	CLI.Flags().StringVarP(&config, "config", "c", "/etc/lastbackend/ingress", "/path/to/config.yml")
	CLI.Flags().StringVarP(&token, "token", "t", "", "Last.Backend cluster authentication token")

	CLI.Flags().BoolVar(&daemon, "daemon", false, "Run hoarder as a server")
	CLI.Flags().IntVarP(&debug, "verbose", "v", 0, "verbose level")

	viper.BindPFlag("verbose", CLI.Flags().Lookup("verbose"))
	viper.BindPFlag("token", CLI.Flags().Lookup("token"))
}

func main() {
	if err := CLI.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
