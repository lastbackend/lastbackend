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
	"strconv"
	"strings"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

func init() {
	serviceUpdateCmd.Flags().StringP("desc", "d", "", "set service description")
	serviceUpdateCmd.Flags().Int64P("memory", "m", 0, "set service spec memory")
	serviceUpdateCmd.Flags().IntP("replicas", "r", 0, "set service replicas")
	serviceUpdateCmd.Flags().StringArrayP("port", "p", make([]string, 0), "set service ports")
	serviceUpdateCmd.Flags().StringArrayP("envs", "e", make([]string, 0), "set service envs")
	serviceUpdateCmd.Flags().StringP("image", "i", "", "set service image")
	serviceUpdateCmd.Flags().StringP("auth", "a", "", "set service image auth")
	serviceCmd.AddCommand(serviceUpdateCmd)
}

const serviceUpdateExample = `
  # Update info for 'redis' service in 'ns-demo' namespace 
  lb service update ns-demo redis --desc "Example new description" -m 128
`

var serviceUpdateCmd = &cobra.Command{
	Use:     "update [NAMESPACE] [NAME]",
	Short:   "Change configuration of the service",
	Example: serviceUpdateExample,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]
		name := args[1]

		description, _ := cmd.Flags().GetString("desc")
		memory, _ := cmd.Flags().GetInt64("memory")
		ports, _ := cmd.Flags().GetStringArray("port")
		env, _ := cmd.Flags().GetStringArray("envs")
		replicas, _ := cmd.Flags().GetInt("replicas")
		image, _ := cmd.Flags().GetString("image")
		secret, _ := cmd.Flags().GetString("auth")

		opts := new(request.ServiceUpdateOptions)
		opts.Spec = new(request.ServiceOptionsSpec)

		if description != "" {
			opts.Description = &description
		}

		if memory != 0 {
			opts.Spec.Memory = &memory
		}

		if replicas != 0 {
			opts.Spec.Replicas = &replicas
		}

		if len(ports) > 0 {
			opts.Spec.Ports = make(map[uint16]string, 0)

			for _, p := range ports {
				pm := strings.Split(p, ":")
				if len(pm) != 2 {
					fmt.Println("port mapping is in invalid format")
					return
				}

				ext, err := strconv.ParseUint(pm[0], 10, 16)
				if err != nil {
					fmt.Println("port mapping is in invalid format")
					return
				}

				opts.Spec.Ports[uint16(ext)] = pm[1]
			}
		}

		if len(env) > 0 {
			opts.Spec.EnvVars = &env
		}

		if image != types.EmptyString {
			if opts.Image == nil {
				opts.Image = new(request.ServiceImageSpec)
			}

			opts.Image.Name = &image
		}

		if secret != types.EmptyString {
			if opts.Image == nil {
				opts.Image = new(request.ServiceImageSpec)
			}
			opts.Image.Secret = &secret
		}

		if err := opts.Validate(); err != nil {
			fmt.Println(err.Err())
			return
		}

		cli := envs.Get().GetClient()
		response, err := cli.V1().Namespace(namespace).Service(name).Update(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(fmt.Sprintf("Service `%s` is updated", name))
		ss := view.FromApiServiceView(response)
		ss.Print()
	},
}
