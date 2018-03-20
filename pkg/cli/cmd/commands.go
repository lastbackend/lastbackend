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
	"github.com/spf13/cobra"

	"fmt"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Client version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(RootCmd.Use + " " + version)
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

var repo = &cobra.Command{
	Use:   "repo",
	Short: "Manage repo(registry)",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
