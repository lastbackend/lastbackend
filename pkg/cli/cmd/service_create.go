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
	serviceCreateCmd.Flags().StringP("desc", "d", "", "set service description")
	serviceCreateCmd.Flags().StringP("name", "n", "", "set service name")
	serviceCreateCmd.Flags().StringP("auth", "a", "", "service image auth secret")
	serviceCreateCmd.Flags().Int64P("memory", "m", 128, "set service spec memory")
	serviceCreateCmd.Flags().IntP("replicas", "r", 1, "set service replicas")
	serviceCreateCmd.Flags().StringArrayP("port", "p", make([]string, 0), "set service ports")
	serviceCreateCmd.Flags().StringArrayP("env", "e", make([]string, 0), "set service env")
	serviceCreateCmd.Flags().StringArray("env-from-secret", make([]string, 0), "set service env from secret")
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
		ports, _ := cmd.Flags().GetStringArray("ports")
		env, _ := cmd.Flags().GetStringArray("env")
		senv, _ := cmd.Flags().GetStringArray("env-from-secret")
		replicas, _ := cmd.Flags().GetInt("replicas")
		auth, _ := cmd.Flags().GetString("auth")

		opts := new(request.ServiceCreateOptions)
		opts.Image = new(request.ServiceImageSpec)
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

		es := make(map[string]request.ServiceEnvOption)
		if len(env) > 0 {
			for _, e := range env {
				kv := strings.SplitN(e, "=", 2)
				eo := request.ServiceEnvOption{
					Name: kv[0],
				}
				if len(kv) > 1 {
					eo.Value = kv[1]
				}

				es[eo.Name] = eo
			}


		}
		if len(senv) > 0 {
			for _, e := range senv {
				kv := strings.SplitN(e, "=", 3)
				eo := request.ServiceEnvOption{
					Name: kv[0],
				}
				if len(kv) < 3 {
					fmt.Println("Service env from secret is in wrong format, should be [NAME]=[SECRET NAME]=[SECRET STORAGE KEY]")
					return
				}

				if len(kv) == 3 {
					eo.From.Name = kv[1]
					eo.From.Key = kv[2]
				}

				es[eo.Name] = eo
			}
		}

		if len(es) > 0 {
			senvs := make([]request.ServiceEnvOption, 0)
			for _, e := range es {
				senvs = append(senvs, e)
			}
			opts.Spec.EnvVars = &senvs
		}


		opts.Description = &description
		opts.Image.Name = &image
		if auth != types.EmptyString {
			opts.Image.Secret = &auth
		}

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
