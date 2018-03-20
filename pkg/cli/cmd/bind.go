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
	ns "github.com/lastbackend/lastbackend/pkg/cli/cmd/namespace"
	sr "github.com/lastbackend/lastbackend/pkg/cli/cmd/service"
)

func init() {
	commands()
	flags()
}

func commands() {

	// ----- root -----
	RootCmd.AddCommand(
		versionCmd,
		namespace,
		service,
		repo,
	)

	// ----- namespace -----
	namespace.AddCommand(
		ns.NamespaceCreate,
		ns.NamespaceList,
		ns.NamespaceInfo,
		ns.NamespaceRemove,
		ns.NamespaceSelect,
	)

	// ----- service -----
	service.AddCommand(
		sr.ServiceList,
		sr.ServiceInfo,
		sr.ServiceScale,
		sr.ServiceUpdate,
		sr.ServiceRemove,
		sr.ServiceCreate,
		sr.ServiceLogs,
	)
}

func flags() {

	// ----- namespace -----
	ns.NamespaceCreate.Flags().StringVarP(&ns.Desc, "desc", "d", "", "Set description")

	// ----- service -----
	sr.ServiceCreate.Flags().Int64VarP(&sr.Memory, "memory", "m", 0, "Set memory")
	sr.ServiceCreate.Flags().StringVarP(&sr.Sources, "sources", "s", "", "Set sources")
	sr.ServiceUpdate.Flags().Int64VarP(&sr.Memory, "memory", "m", 0, "Set memory")
}
