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
	"strings"
	"strconv"
)

func init() {
	routeCmd.AddCommand(routeCreateCmd)
}

const routeCreateExample = `
  # Create new route for proxying http traffic from 'blog-ns-demo.lstbknd.io' to service 'blog-web' on 80 port
  lb route create ns-demo blog blog-web:80"
`

var routeCreateCmd = &cobra.Command{
	Use:     "create [NAMESPACE] [NAME] [SERVICE:PORT]",
	Short:   "Create new route",
	Example: routeCreateExample,
	Args:    cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]

		opts := new(request.RouteCreateOptions)
		opts.Name = args[1]

		proxy := strings.Split(args[2], ":")
		port, err := strconv.Atoi(proxy[1])
		if err != nil {
			fmt.Printf("Invalid port number: %s", proxy[1])
			return
		}

		if port >= 65535 {
			fmt.Printf("Port number is out of range: %s [65535]", proxy[1])
			return
		}

		opts.Rules = append(opts.Rules, request.RulesOption{
			Service: proxy[0],
			Port: port,
		})


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

		fmt.Println(fmt.Sprintf("Route `%s` is created in namespace `%s`", opts.Name, namespace))

		service := view.FromApiRouteView(response)
		service.Print()
	},
}
