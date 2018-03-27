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
// patents in process, and are protected by trade secretCmd or copyright law.
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
	"github.com/lastbackend/lastbackend/pkg/cli/envs"
	"github.com/lastbackend/lastbackend/pkg/cli/storage"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(
		clusterCmd,
		namespaceCmd,
		routeCmd,
		serviceCmd,
		secretCmd,
		tokenCmd,
		versionCmd,
	)
}

var (
	cfg = config.Get()
	ctx = envs.Get()
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "lb",
	Short: "Apps cloud hosting with integrated deployment tools",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		debug, _ := cmd.Flags().GetBool("debug")

		if debug {
			fmt.Println("debug mode enabled")
			cfg.Debug = debug
		}

		token, err := storage.GetToken()
		if err != nil {
			panic("There is no token in .lastbackend in homedir")
		}

		host := cmd.Flag("host").Value.String()

		httpcli, err := client.NewHTTP(host, &client.Config{BearerToken: token})
		if err != nil {
			panic(err)
		}

		ctx.SetClient(httpcli)
	},
}

var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Manage your namespaces",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		command := "[string]"
		if len(args) != 0 {
			command = args[0]
		}

		var ns = &cobra.Command{
			Use:   command,
			Short: "Manage your a namespace",
		}

		cmd.AddCommand(ns)

		if len(args) == 0 {
			cmd.Help()
			return
		}

		// Attach sub command for namespace
		ns.AddCommand(
			serviceCmd,
			secretCmd,
			routeCmd,
		)

		ns.Execute()

	},
}

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage your service",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage your secret",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var routeCmd = &cobra.Command{
	Use:   "route",
	Short: "Manage your route",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage your cluster",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage set vars to your local storage",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command
func Execute() {

	cobra.OnInitialize()

	RootCmd.PersistentFlags().StringP("host", "H", "https://api.lastbackend.com", "Set api host parameter")
	RootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
