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

var Desc string

var NamespaceCreate = &cobra.Command{
	Use:   "create",
	Short: "Create new namespace",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Help()
			return
		}

		CreateCmd(args[0], Desc)
	},
}

var NamespaceList = &cobra.Command{
	Use:   "list",
	Short: "Display the Namespace list",
	Run: func(cmd *cobra.Command, args []string) {
		ListCmd()
	},
}

var NamespaceFetch = &cobra.Command{
	Use:   "info",
	Short: "Get namespace info by Name, if without Name - get current namespace info",
	Run: func(cmd *cobra.Command, args []string) {
		FetchCmd(args[0])
	},
}

var NamespaceRemove = &cobra.Command{
	Use:   "remove",
	Short: "Remove namespace by Name",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Help()
			return
		}

		RemoveCmd(args[0])
	},
}

var NamespaceCurrent = &cobra.Command{
	Use:   "current",
	Short: "Show current namespace",
	Run: func(cmd *cobra.Command, args []string) {
		CurrentCmd()
	},
}

var NamespaceSelect = &cobra.Command{
	Use:   "select",
	Short: "Select to the namespace",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Help()
			return
		}

		SelectCmd(args[0])
	},
}
