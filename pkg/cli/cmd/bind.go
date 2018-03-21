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
	cs "github.com/lastbackend/lastbackend/pkg/cli/cmd/cluster"
	ns "github.com/lastbackend/lastbackend/pkg/cli/cmd/namespace"
	sr "github.com/lastbackend/lastbackend/pkg/cli/cmd/service"
	st "github.com/lastbackend/lastbackend/pkg/cli/cmd/set"
)

func init() {
	// ----- root -----
	RootCmd.AddCommand(
		versionCmd,
		cluster,
		namespace,
		service,
	)

	// ----- cluster -----
	cluster.AddCommand(
		cs.ClusterFetch,
	)

	// ----- set -----
	set.AddCommand(
		st.SetToken,
	)

	// ----- namespace -----
	namespace.AddCommand(
		ns.NamespaceCreate,
		ns.NamespaceList,
		ns.NamespaceFetch,
		ns.NamespaceRemove,
	)

	// ----- service -----
	service.AddCommand(
		sr.ServiceCreate,
		sr.ServiceList,
		sr.ServiceFetch,
		sr.ServiceUpdate,
		sr.ServiceRemove,
	)
}
