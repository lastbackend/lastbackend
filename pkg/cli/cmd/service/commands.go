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

var ServiceCreate = &cobra.Command{
	Use:   "create",
	Short: "Create service",
	Run:   CreateCmd,
}

var ServiceFetch = &cobra.Command{
	Use:   "info",
	Short: "Service info by name",
	Run:   FetchCmd,
}

var ServiceList = &cobra.Command{
	Use:   "list",
	Short: "Display the services list",
	Run:   ListCmd,
}

var ServiceUpdate = &cobra.Command{
	Use:   "update",
	Short: "Change configuration of the service",
	Run:   UpdateCmd,
}

var ServiceRemove = &cobra.Command{
	Use:   "remove",
	Short: "Remove service by Name",
	Run:   RemoveCmd,
}
