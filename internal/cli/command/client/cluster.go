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
	"errors"
	"fmt"

	"github.com/lastbackend/lastbackend/internal/cli/service"
	"github.com/lastbackend/lastbackend/internal/cli/views"
	"github.com/lastbackend/lastbackend/tools/logger"
	"github.com/spf13/cobra"
)

const clusterAddExample = `
  # Get information about cluster 
  lb cluster add name endpoint --local
`

const clusterDelExample = `
  # Get information about cluster 
  lb cluster del name --local
`

const clusterInspectExample = `
  # Get information about cluster 
  lb cluster inspect
`

const clusterListExample = `
  # Get information about cluster 
  lb cluster ls
`

const clusterSelectExample = `
  # Get information about cluster 
  lb cluster select name
`

func (c *command) NewClusterCmd(clusterService *service.ClusterService) *cobra.Command {
	log := logger.WithContext(context.Background())

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Manage your cluster",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Error(err.Error())
				return
			}
		},
	}

	cmd.AddCommand(c.clusterAddCmd(clusterService))
	cmd.AddCommand(c.clusterDelCmd(clusterService))
	cmd.AddCommand(c.clusterInspectCmd(clusterService))
	cmd.AddCommand(c.clusterListCmd(clusterService))
	cmd.AddCommand(c.clusterSelectCmd(clusterService))

	return cmd
}

func (c *command) clusterAddCmd(clusterService *service.ClusterService) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [NAME] [ENDPOINT]",
		Short:   "Add cluster",
		Example: clusterAddExample,
		Args:    cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			local, _ := cmd.Flags().GetBool("local")
			if !local {
				return errors.New("method allowed with local flag")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			name := args[0]
			endpoint := args[1]

			local, err := cmd.Flags().GetBool("local")
			if err != nil {
				panic(err)
			}

			token, err := cmd.Flags().GetString("token")
			if err != nil {
				panic(err)
			}

			err = clusterService.AddLocalCluster(name, endpoint, token, local)
			switch true {
			case err == nil:
			case err.Error() == "already exists":
				fmt.Println(fmt.Sprintf("Cluster `%s` already exists", name))
			default:
				panic(err)
			}

		},
	}

	cmd.Flags().Bool("local", false, "Use local cluster")

	return cmd
}

func (c *command) clusterDelCmd(clusterService *service.ClusterService) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "del [NAME]",
		Short:   "Remove cluster",
		Example: clusterDelExample,
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			local, _ := cmd.Flags().GetBool("local")
			if !local {
				return errors.New("method allowed with local flag")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			if err := clusterService.DelLocalCluster(name); err != nil {
				panic(err)
			}

		},
	}

	cmd.Flags().Bool("local", false, "Use local cluster")

	return cmd
}

func (c *command) clusterInspectCmd(clusterService *service.ClusterService) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inspect",
		Short:   "Get cluster info",
		Example: clusterInspectExample,
		Args:    cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {

			response, err := c.client.cluster.V1().Cluster().Get(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			views.FromLbApiClusterView(response).Print()
		},
	}

	return cmd
}

func (c *command) clusterListCmd(clusterService *service.ClusterService) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "Get available cluster list",
		Example: clusterListExample,
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println("Remotely clusters list:")

			ritems, err := c.client.genesis.V1().Cluster().List(context.Background())
			if err != nil {
				fmt.Println(err)
			}

			vg := views.FromGenesisApiClusterListView(ritems)
			vg.Print()

			fmt.Print("\n")
			fmt.Println("Locally clusters list:")

			litems, err := clusterService.List()
			if err != nil {
				fmt.Println(err)
				return
			}

			vs := views.FromStorageClusterList(litems)
			vs.Print()
		},
	}

	return cmd
}

func (c *command) clusterSelectCmd(clusterService *service.ClusterService) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "select [NAME]",
		Short:   "Select cluster",
		Example: clusterSelectExample,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cmd.Help()
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			name := args[0]

			// Select local cluster
			local, _ := cmd.Flags().GetBool("local")
			if local {

				item, err := clusterService.GetLocalCluster(name)
				if err != nil {
					fmt.Println(err)
				}

				if item == nil {
					fmt.Println(fmt.Sprintf("Cluster `%s` not found", name))
				}

				err = clusterService.SetCluster(fmt.Sprintf("l.%s", item.Name))
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println(fmt.Sprintf("Cluster `%s` selected", name))

				return
			}

			cl, err := c.client.genesis.V1().Cluster().Get(context.Background(), name)
			if err != nil {
				fmt.Println(err)
				return
			}

			if cl == nil {
				fmt.Println(fmt.Sprintf("Cluster `%s` not found", name))
				return
			}

			err = clusterService.SetCluster(fmt.Sprintf("r.%s", cl.Meta.SelfLink))
			if err != nil {
				panic(err)
			}

			fmt.Println(fmt.Sprintf("Cluster `%s` selected", name))
		},
	}

	cmd.Flags().Bool("local", false, "Use local cluster")

	return cmd
}
