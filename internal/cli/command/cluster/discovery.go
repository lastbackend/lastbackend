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

const discoveryInspectExample = `
  # Get information 'wef34fg' for discovery
  lb discovery inspect wef34fg"
`

const discoveryListExample = `
  # Get all discoveryes in cluster 
  lb discovery ls
`

func (c *command) NewDiscoveryCmd() *cobra.Command {
	log := logger.WithContext(context.Background())
	cmd := &cobra.Command{
		Use:   "discovery",
		Short: "Manage cluster discovery servers",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Error(err.Error())
				return
			}
		},
	}

	cmd.AddCommand(c.discoveryInspectCmd())
	cmd.AddCommand(c.discoveryListCmd())

	return cmd
}

func (c *command) discoveryInspectCmd() *cobra.Command {
	cmd:= &cobra.Command{
		Use:     "inspect [NAME]",
		Short:   "Discovery info by name",
		Example: discoveryInspectExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			name := args[0]

			response, err := c.client.cluster.V1().Cluster().Discovery(name).Get(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			ss := views.FromApiDiscoveryView(response)
			ss.Print()
		},
	}

	return cmd
}

func (c *command) discoveryListCmd() *cobra.Command {
	cmd:= &cobra.Command{
		Use:     "ls",
		Short:   "Display the discoveries list",
		Example: discoveryListExample,
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			response, err := c.client.cluster.V1().Cluster().Discovery().List(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			if response == nil || len(*response) == 0 {
				fmt.Println("no discoveries available")
				return
			}

			list := views.FromApiDiscoveryListView(response)
			list.Print()
		},
	}

	return cmd
}
