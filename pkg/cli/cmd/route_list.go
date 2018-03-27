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
	routeCmd.AddCommand(routeListCmd)
}

const routeListExample = `
  # Get all routes for 'ns-demo' namespace
  lb route inspect ns-demo wef34fg"
`

var routeListCmd = &cobra.Command{
	Use:     "ls [NAMESPACE]",
	Short:   "Get routes list",
	Example: routeListExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]

		cli := envs.Get().GetClient()
		response, err := cli.V1().Namespace(namespace).Route().List(envs.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		if response == nil || len(*response) == 0 {
			fmt.Println("no routes available")
			return
		}

		list := view.FromApiRouteListView(response)
		list.Print()
	},
}
