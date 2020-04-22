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
	"github.com/lastbackend/lastbackend/tools/logger"

	"github.com/lastbackend/lastbackend/internal/cli/views"
	"github.com/spf13/cobra"
)

const nodeListExample = `
  # Get all nodes for 'ns-demo' namespace  
  lb node ls
`

const nodeInspectExample = `
  # Get information 'wef34fg' for node in 'ns-demo' namespace
  lb node inspect ns-demo wef34fg"
`

func (c *command) NewNodeCmd() *cobra.Command {
	log := logger.WithContext(context.Background())
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Manage cluster nodes",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Error(err.Error())
				return
			}
		},
	}

	cmd.AddCommand(c.nodeListCmd())
	cmd.AddCommand(c.nodeInspectCmd())

	return cmd
}

func (c *command) nodeListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "Display the nodes list",
		Example: nodeListExample,
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			response, err := c.client.cluster.V1().Cluster().Node().List(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			if response == nil || len(*response) == 0 {
				fmt.Println("no nodes available")
				return
			}

			list := views.FromApiNodeListView(response)
			list.Print()
		},
	}

	return cmd
}

func (c *command) nodeInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inspect [NAME]",
		Short:   "Node info by name",
		Example: nodeInspectExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			name := args[0]

			response, err := c.client.cluster.V1().Cluster().Node(name).Get(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			ss := views.FromApiNodeView(response)
			ss.Print()
		},
	}

	return cmd
}
