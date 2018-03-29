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
// patents in process, and are protected by trade secretCmd or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package cmd

import (
	"fmt"

	"github.com/lastbackend/lastbackend/pkg/cli/envs"
	"github.com/lastbackend/lastbackend/pkg/cli/view"
	"github.com/spf13/cobra"
)

func init() {
	routeCmd.AddCommand(routeInspectCmd)
}

const routeInspectExample = `
  # Get information 'wef34fg' for route in 'ns-demo' namespace
  lb route inspect ns-demo wef34fg"
`

var routeInspectCmd = &cobra.Command{
	Use:     "inspect [NAMESPACE] [NAME]",
	Short:   "Route info by name",
	Example: routeInspectExample,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]
		name := args[1]

		cli := envs.Get().GetClient()
		response, err := cli.V1().Namespace(namespace).Route(name).Get(envs.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		ss := view.FromApiRouteView(response)
		ss.Print()
	},
}
