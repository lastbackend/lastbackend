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

package service

import (
	"github.com/spf13/cobra"
)

var (
	Sources string
	Memory  int64
)

var ServiceList = &cobra.Command{
	Use:   "list",
	Short: "Display the services list",
	Run: func(cmd *cobra.Command, args []string) {
		//ListServiceCmd()
	},
}

var ServiceScale = &cobra.Command{
	Use:   "scale",
	Short: "Scale service",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 2 {
			cmd.Help()
			return
		}

		//replicas, err := strconv.ParseInt(args[1], 10, 64)
		//if err != nil {
		//	cmd.Help()
		//	return
		//}

		//ScaleCmd(args[0], replicas)
	},
}

var ServiceInfo = &cobra.Command{
	Use:   "info",
	Short: "Service info by Name",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Help()
			return
		}

		//InspectCmd(args[0])
	},
}

var ServiceUpdate = &cobra.Command{
	Use:   "update",
	Short: "Change configuration of the service",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Help()
			return
		}

		//UpdateCmd(args[0], Memory)
	},
}

var ServiceRemove = &cobra.Command{
	Use:   "remove",
	Short: "Remove service by Name",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Help()
			return
		}

		//RemoveCmd(args[0])
	},
}

var ServiceCreate = &cobra.Command{
	Use:   "create",
	Short: "Create service",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Help()
			return
		}

		//CreateCmd(args[0], Sources, Memory)
	},
}

var ServiceLogs = &cobra.Command{
	Use:   "logs",
	Short: "Show service logs",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Help()
			return
		}

		//LogsServiceCmd(args[0])
	},
}
