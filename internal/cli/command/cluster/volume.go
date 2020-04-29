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
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/spf13/cobra"
)

const volumeListExample = `
  # Get all volumes for 'ns-demo' namespace  
  lb volume ls ns-demo
`

const volumeInspectExample = `
  # Get information for 'redis' volume in 'ns-demo' namespace
  lb volume inspect ns-demo redis
`

const volumeCreateExample = `
  # Create new redis volume with description and 256 MB limit memory
  lb volume create ns-demo redis --desc "Example description" -m 256
`

const volumeRemoveExample = `
  # Remove 'redis' volume in 'ns-demo' namespace
  lb volume remove ns-demo redis
`

const volumeUpdateExample = `
  # Update info for 'redis' volume in 'ns-demo' namespace
  lb volume update ns-demo redis --desc "Example new description" -m 128
`

func (c *command) NewVolumeCmd() *cobra.Command {
	log := logger.WithContext(context.Background())
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "Manage your volumes",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Error(err.Error())
				return
			}
		},
	}

	cmd.AddCommand(c.volumeListCmd())
	cmd.AddCommand(c.volumeInspectCmd())
	cmd.AddCommand(c.volumeCreateCmd())
	cmd.AddCommand(c.volumeRemoveCmd())
	cmd.AddCommand(c.volumeUpdateCmd())

	return cmd
}

func (c *command) volumeListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls [NAMESPACE]",
		Short:   "Display the volumes list",
		Example: volumeListExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]

			response, err := c.client.cluster.V1().Namespace(namespace).Volume().List(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			if response == nil || len(*response) == 0 {
				fmt.Println("no volumes available")
				return
			}

			list := views.FromApiVolumeListView(response)
			list.Print()
		},
	}
}

func (c *command) volumeInspectCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "inspect [NAMESPACE] [NAME]",
		Short:   "Volume info by name",
		Example: volumeInspectExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]
			name := args[1]

			response, err := c.client.cluster.V1().Namespace(namespace).Volume(name).Get(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			ss := views.FromApiVolumeView(response)
			ss.Print()
		},
	}
}

func (c *command) volumeCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "create [NAMESPACE] [NAME]",
		Short:   "Create volume",
		Example: volumeCreateExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]
			name := args[1]

			description, _ := cmd.Flags().GetString("desc")
			kind, _ := cmd.Flags().GetString("type")

			opts := new(request.VolumeManifest)

			if len(name) != 0 {
				opts.Meta.Name = &name
			}

			if len(description) != 0 {
				opts.Meta.Description = &description
			}

			switch kind {
			case models.KindVolumeHostDir:
				opts.Spec.Type = models.KindVolumeHostDir
				break
			}

			if err := opts.Validate(); err != nil {
				fmt.Println(err.Err())
				return
			}

			response, err := c.client.cluster.V1().Namespace(namespace).Volume().Create(context.Background(), opts)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("Volume `%s` is created", name))

			volume := views.FromApiVolumeView(response)
			volume.Print()
		},
	}
}

func (c *command) volumeRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove [NAMESPACE] [NAME]",
		Short:   "Remove volume by name",
		Example: volumeRemoveExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]
			name := args[1]

			opts := &request.VolumeRemoveOptions{Force: false}

			if err := opts.Validate(); err != nil {
				fmt.Println(err.Err())
				return
			}

			c.client.cluster.V1().Namespace(namespace).Volume(name).Remove(context.Background(), opts)

			fmt.Println(fmt.Sprintf("Volume `%s` is successfully removed", name))
		},
	}
}

func (c *command) volumeUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "update [NAMESPACE] [NAME]",
		Short:   "Change configuration of the volume",
		Example: volumeUpdateExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]
			name := args[1]

			description, _ := cmd.Flags().GetString("desc")

			opts := new(request.VolumeManifest)

			if len(name) != 0 {
				opts.Meta.Name = &name
			}

			if len(description) != 0 {
				opts.Meta.Description = &description
			}

			if err := opts.Validate(); err != nil {
				fmt.Println(err.Err())
				return
			}

			response, err := c.client.cluster.V1().Namespace(namespace).Volume(name).Update(context.Background(), opts)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("Volume `%s` is updated", name))
			ss := views.FromApiVolumeView(response)
			ss.Print()
		},
	}
}
