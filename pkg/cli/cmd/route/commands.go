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

package route

import (
	"github.com/spf13/cobra"
)

var RouteCreate = &cobra.Command{
	Use:   "create",
	Short: "Create route",
	Run:   CreateCmd,
}

var RouteFetch = &cobra.Command{
	Use:   "info",
	Short: "Route info by name",
	Run:   FetchCmd,
}

var RouteList = &cobra.Command{
	Use:   "list",
	Short: "Display the routes list",
	Run:   ListCmd,
}

var RouteUpdate = &cobra.Command{
	Use:   "update",
	Short: "Change configuration of the route",
	Run:   UpdateCmd,
}

var RouteRemove = &cobra.Command{
	Use:   "remove",
	Short: "Remove route by Name",
	Run:   RemoveCmd,
}
