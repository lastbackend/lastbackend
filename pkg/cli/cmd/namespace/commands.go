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

var NamespaceCreate = &cobra.Command{
	Use:   "create",
	Short: "Create new namespace",
	Run:   CreateCmd,
}

var NamespaceFetch = &cobra.Command{
	Use:   "info",
	Short: "Get namespace info by name",
	Run:   FetchCmd,
}

var NamespaceList = &cobra.Command{
	Use:   "list",
	Short: "Display the namespace list",
	Run:   ListCmd,
}

var NamespaceUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update the namespace by name",
	Run:   UpdateCmd,
}

var NamespaceRemove = &cobra.Command{
	Use:   "remove",
	Short: "Remove namespace by name",
	Run:   RemoveCmd,
}
