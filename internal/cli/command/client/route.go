//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package client

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/tools/logger"
	"strconv"
	"strings"

	"github.com/lastbackend/lastbackend/internal/cli/views"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/spf13/cobra"
)

const routeListExample = `
  # Get all routes for 'ns-demo' namespace
  lb route inspect ns-demo wef34fg"
`

const routeInspectExample = `
  # Get information 'wef34fg' for route in 'ns-demo' namespace
  lb route inspect ns-demo wef34fg"
`

const routeCreateExample = `
  # Create new route for proxying http traffic from 'blog-ns-demo.lstbknd.io' to service 'blog-web' on 80 port
  lb route create ns-demo blog blog-web:80"
`

const routeRemoveExample = `
  # Remove 'wef34fg' route for 'ns-demo' namespace
  lb route remove ns-demo wef34fg"
`

const routeUpdateExample = `
  # Update 'wef34fg' route for 'ns-demo' namespace
  lb route update ns-demo wef34fg blog-ns-demo.lstbknd.net 443"
`

func (c *command) NewRouteCmd() *cobra.Command {
	log := logger.WithContext(context.Background())
	cmd := &cobra.Command{
		Use:   "route",
		Short: "Manage your route",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Error(err.Error())
				return
			}
		},
	}

	cmd.AddCommand(c.routeListCmd())
	cmd.AddCommand(c.routeInspectCmd())
	cmd.AddCommand(c.routeCreateCmd())
	cmd.AddCommand(c.routeRemoveCmd())
	cmd.AddCommand(c.routeUpdateCmd())

	return cmd
}

func (c *command) routeListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls [NAMESPACE]",
		Short:   "Get routes list",
		Example: routeListExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]

			response, err := c.client.cluster.V1().Namespace(namespace).Route().List(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			if response == nil || len(*response) == 0 {
				fmt.Println("no routes available")
				return
			}

			list := views.FromApiRouteListView(response)
			list.Print()
		},
	}

	return cmd
}

func (c *command) routeInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inspect [NAMESPACE] [NAME]",
		Short:   "Route info by name",
		Example: routeInspectExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]
			name := args[1]

			response, err := c.client.cluster.V1().Namespace(namespace).Route(name).Get(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			ss := views.FromApiRouteView(response)
			ss.Print()
		},
	}

	return cmd
}

func (c *command) routeCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [NAMESPACE] [NAME] [SERVICE:PORT]",
		Short:   "Create new route",
		Example: routeCreateExample,
		Args:    cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]

			opts := new(request.RouteManifest)
			opts.Meta.Name = &args[1]

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

			opts.Spec.Rules = append(opts.Spec.Rules, request.RouteManifestSpecRulesOption{
				Service: proxy[0],
				Port:    port,
			})

			if err := opts.Validate(); err != nil {
				fmt.Println(err.Err())
				return
			}

			response, err := c.client.cluster.V1().Namespace(namespace).Route().Create(context.Background(), opts)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("Route `%s` is created in namespace `%s`", *opts.Meta.Name, namespace))

			service := views.FromApiRouteView(response)
			service.Print()
		},
	}

	return cmd
}

func (c *command) routeRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove [NAMESPACE] [NAME]",
		Short:   "Remove route by name",
		Example: routeRemoveExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]
			name := args[1]

			opts := &request.RouteRemoveOptions{Force: false}

			if err := opts.Validate(); err != nil {
				fmt.Println(err.Err())
				return
			}

			c.client.cluster.V1().Namespace(namespace).Route(name).Remove(context.Background(), opts)

			fmt.Println(fmt.Sprintf("Route `%s` remove now", name))
		},
	}

	return cmd
}

func (c *command) routeUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update [NAMESPACE] [NAME] [ENDPOINT] [PORT]",
		Short:   "Change configuration of the route",
		Example: routeUpdateExample,
		Args:    cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]
			name := args[1]

			// TODO: set routeCmd options
			opts := new(request.RouteManifest)

			if err := opts.Validate(); err != nil {
				fmt.Println(err.Err())
				return
			}

			response, err := c.client.cluster.V1().Namespace(namespace).Route(name).Update(context.Background(), opts)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("Route `%s` is updated", name))
			ss := views.FromApiRouteView(response)
			ss.Print()
		},
	}

	return cmd
}
