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

	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/cli/envs"
	"github.com/lastbackend/lastbackend/pkg/cli/view"
	"github.com/spf13/cobra"
)

func init() {
	routeCmd.AddCommand(routeCreateCmd)
}

const routeCreateExample = `
  # Create new route for proxing http traffic from 'blog-ns-demo.lstbknd.net' to service on 80 port
  lb route create ns-demo blog-ns-demo.lstbknd.net 80"
`

var routeCreateCmd = &cobra.Command{
	Use:     "create [NAMESPACE] [ENDPOINT] [PORT]",
	Short:   "Create new route",
	Example: routeCreateExample,
	Args:    cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]
		name := args[1]

		// TODO: set route options

		opts := new(request.RouteCreateOptions)

		if err := opts.Validate(); err != nil {
			fmt.Println(err.Err())
			return
		}

		cli := envs.Get().GetClient()
		response, err := cli.V1().Namespace(namespace).Route().Create(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(fmt.Sprintf("Route `%s` is created", name))

		service := view.FromApiRouteView(response)
		service.Print()
	},
}
