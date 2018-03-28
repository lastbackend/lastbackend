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
	serviceCreateCmd.Flags().StringP("desc", "d", "", "set service description")
	serviceCreateCmd.Flags().StringP("name", "n", "", "set service name")
	serviceCreateCmd.Flags().Int64P("memory", "m", 128, "set service spec memory")
	serviceCreateCmd.Flags().IntP("replicas", "r", 1, "set service replicas")
	serviceCmd.AddCommand(serviceCreateCmd)
}

const serviceCreateExample = `
  # Create new redis service with description and 256 MB limit memory
  lb service create ns-demo redis --desc "Example description" -m 256
`

var serviceCreateCmd = &cobra.Command{
	Use:     "create [NAMESPACE] [IMAGE]",
	Short:   "Create service",
	Example: serviceCreateExample,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]
		image := args[1]

		description, _ := cmd.Flags().GetString("desc")
		memory, _ := cmd.Flags().GetInt64("memory")
		name, _ := cmd.Flags().GetString("name")
		replicas, _ := cmd.Flags().GetInt("replicas")

		opts := new(request.ServiceCreateOptions)
		opts.Spec = new(request.ServiceOptionsSpec)

		if len(name) != 0 {
			opts.Name = &name
		}

		if len(description) != 0 {
			opts.Description = &description
		}

		if memory != 0 {
			opts.Spec.Memory = &memory
		}

		if replicas != 0 {
			opts.Spec.Replicas = &replicas
		}

		opts.Description = &description
		opts.Image = &image

		if err := opts.Validate(); err != nil {
			fmt.Println(err.Err())
			return
		}

		cli := envs.Get().GetClient()
		response, err := cli.V1().Namespace(namespace).Service().Create(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(fmt.Sprintf("Service `%s` is created", name))

		service := view.FromApiServiceView(response)
		service.Print()
	},
}
