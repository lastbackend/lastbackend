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
	commands()
	flags()
}

func commands() {
	// ----- root -----
	RootCmd.AddCommand(
		versionCmd,
		cluster,
		namespace,
		service,
	)

	RootCmd.Flags().BoolP("dd", "a", true, "a")

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
		ns.NamespaceUpdate,
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

func flags() {
	// ----- NAMESPACE -----
	// namespace :: NamespaceCreate
	ns.NamespaceCreate.Flags().StringP("desc", "d", "", "set namespace description")
	// namespace :: NamespaceUpdate
	ns.NamespaceUpdate.Flags().StringP("desc", "d", "", "set namespace description")

	// ----- SERVICE -----
	// service :: ServiceCreate
	sr.ServiceCreate.Flags().String("namespace", "", "set namespace context")
	sr.ServiceCreate.Flags().StringP("desc", "d", "", "dset service description")
	sr.ServiceCreate.Flags().StringP("image", "i", "", "set service spec image")
	sr.ServiceCreate.Flags().Int64P("memory", "m", 128, "set service spec memory")
	sr.ServiceCreate.Flags().IntP("replicas", "r", 1, "set service replicas")
	//// service :: ServiceFetch
	sr.ServiceFetch.Flags().String("namespace", "", "set namespace context")
	//// service :: ServiceList
	sr.ServiceList.Flags().String("namespace", "", "set namespace context")
	//// service :: ServiceUpdate
	sr.ServiceUpdate.Flags().String("namespace", "", "set namespace context")
	sr.ServiceUpdate.Flags().StringP("desc", "d", "", "set service description")
	sr.ServiceUpdate.Flags().Int64P("memory", "m", 128, "set service spec memory")
	//// service :: ServiceRemove
	sr.ServiceRemove.Flags().String("namespace", "", "set namespace context")

}
