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

package cmd

import (
	"fmt"
	"os"

	"github.com/lastbackend/lastbackend/pkg/api/client"
	"github.com/lastbackend/lastbackend/pkg/cli/config"
	"github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/spf13/cobra"
)

var (
	// VERSION is set during build
	version string
	host    string
	debug   bool
	tls     bool
)
var (
	cfg = config.Get()
	ctx = context.Get()
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "lb",
	Short: "Apps cloud hosting with integrated deployment tools",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		if debug {
			cfg.Debug = debug
		}
	},
}

// Execute adds all child commands to the root command
func Execute() {
	version = "0.0.1"

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// init client object
func InitClient(apiHost string) *client.Client {
	cli, err := client.New(apiHost)
	if err != nil {
		panic(err)
	}
	return cli
}

func init() {
	cobra.OnInitialize()

	RootCmd.Flags().StringVar(&host, "host", "https://api.lastbackend.com", "Set api host parameter")
	RootCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug mode")
	RootCmd.Flags().BoolVar(&tls, "tls", false, "Enable tls")

	ctx.SetClient(InitClient("https://api.lstbknd.net"))
}
