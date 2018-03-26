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
	serviceCreateCmd.Flags().StringP("desc", "d", "", "dset serviceCmd description")
	serviceCreateCmd.Flags().StringP("name", "n", "", "set serviceCmd name")
	serviceCreateCmd.Flags().Int64P("memory", "m", 128, "set serviceCmd spec memory")
	serviceCreateCmd.Flags().IntP("replicas", "r", 1, "set serviceCmd replicas")
	serviceCmd.AddCommand(serviceCreateCmd)
}

var serviceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create serviceCmd",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Help()
			return
		}

		namespace := cmd.Parent().Parent().Name()

		if namespace == "" {
			fmt.Println("namesapace parameter not set")
			return
		}

		description, _ := cmd.Flags().GetString("desc")
		memory, _ := cmd.Flags().GetInt64("memory")
		name, _ := cmd.Flags().GetString("name")
		replicas, _ := cmd.Flags().GetInt("replicas")

		opts := new(request.ServiceCreateOptions)
		opts.Spec = new(request.ServiceOptionsSpec)

		opts.Name = &name
		opts.Description = &description
		opts.Image = &args[0]
		opts.Spec.Memory = &memory
		opts.Spec.Replicas = &replicas

		if err := opts.Validate(); err != nil {
			fmt.Println(err.Attr)
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
