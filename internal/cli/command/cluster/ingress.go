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

package cluster

import (
	"context"
	"fmt"

	"github.com/lastbackend/lastbackend/internal/cli/views"
	"github.com/lastbackend/lastbackend/tools/log"
	"github.com/spf13/cobra"
)

const ingressInspectExample = `
  # Get information 'wef34fg' for ingress
  lb ingress inspect wef34fg"
`

const ingressListExample = `
  # Get all ingresses in cluster 
  lb ingress ls
`

func (c *command) NewIngressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingress",
		Short: "Manage cluster ingress servers",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Error(err.Error())
				return
			}
		},
	}

	cmd.AddCommand(c.ingressInspectCmd())
	cmd.AddCommand(c.ingressListCmd())

	return cmd
}

func (c *command) ingressInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inspect [NAME]",
		Short:   "Ingress info by name",
		Example: ingressInspectExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			name := args[0]

			response, err := c.client.cluster.V1().Cluster().Ingress(name).Get(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			ss := views.FromApiIngressView(response)
			ss.Print()
		},
	}

	return cmd
}

func (c *command) ingressListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "Display the ingress list",
		Example: ingressListExample,
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			response, err := c.client.cluster.V1().Cluster().Ingress().List(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			if response == nil || len(*response) == 0 {
				fmt.Println("no ingress available")
				return
			}

			list := views.FromApiIngressListView(response)
			list.Print()
		},
	}

	return cmd
}
