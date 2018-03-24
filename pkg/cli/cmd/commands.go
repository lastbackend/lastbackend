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

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Client version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(RootCmd.Use + " " + cmd.Flag("version").Value.String())
	},
}

var set = &cobra.Command{
	Use:   "set",
	Short: "Manage set vars to your local storage",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var cluster = &cobra.Command{
	Use:   "cluster",
	Short: "Manage your cluster",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var namespace = &cobra.Command{
	Use:   "namespace",
	Short: "Manage your namespace and create",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var service = &cobra.Command{
	Use:   "service",
	Short: "Manage service",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var route = &cobra.Command{
	Use:   "route",
	Short: "Manage route ",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var secret = &cobra.Command{
	Use:   "secret",
	Short: "Manage secret ",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
