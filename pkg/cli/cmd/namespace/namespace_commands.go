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

package namespace

import (
	"github.com/spf13/cobra"
)

var (
	Desc string
)

var CreateWorkspace = &cobra.Command{
	Use:   "create",
	Short: "Create new Workspace",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Help()
			return
		}

		CreateCmd(args[0], Desc)
	},
}

var ListWorkspace = &cobra.Command{
	Use:   "list",
	Short: "Display the Workspace list",
	Run: func(cmd *cobra.Command, args []string) {
		ListCmd()
	},
}

var InfoWorkspace = &cobra.Command{
	Use:   "info",
	Short: "Get Workspace info by Name, if without Name - get current Workspace info",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			CurrentCmd("")
		} else {
			CurrentCmd(args[0])
		}
	},
}

var RemoveWorkspace = &cobra.Command{
	Use:   "remove",
	Short: "Remove Workspace by Name",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Help()
			return
		}

		RemoveCmd(args[0])
	},
}

var SelectWorkspace = &cobra.Command{
	Use:   "select",
	Short: "Select to the Workspace",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Help()
			return
		}

		SelectCmd(args[0])
	},
}
