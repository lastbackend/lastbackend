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
	serviceUpdateCmd.Flags().StringArrayP("env", "e", make([]string, 0), "set service env")
	serviceUpdateCmd.Flags().StringArray("env-from-secret", make([]string, 0), "set service env from secret")
	serviceUpdateCmd.Flags().StringP("image", "i", "", "set service image")
	serviceUpdateCmd.Flags().String("image-secret", "", "set service image auth")
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
		env, _ := cmd.Flags().GetStringArray("env")
		senv, _ := cmd.Flags().GetStringArray("env-from-secret")
		replicas, _ := cmd.Flags().GetInt("replicas")
		image, _ := cmd.Flags().GetString("image")
		secret, _ := cmd.Flags().GetString("image-secret")

		opts := new(request.ServiceManifest)
		css := make([]request.ManifestSpecTemplateContainer, 0)
		cs := request.ManifestSpecTemplateContainer{}


		if len(name) != 0 {
			opts.Meta.Name = &name
		}

		if len(description) != 0 {
			opts.Meta.Description = &description
		}


		if memory != 0 {
			cs.Resources.Request.RAM = memory
		}

		if replicas != 0 {
			opts.Spec.Replicas = &replicas
		}

		if len(ports) > 0 {
			opts.Spec.Network = new(request.ManifestSpecNetwork)
			opts.Spec.Network.Ports = make(map[uint16]string, 0)

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

				opts.Spec.Network.Ports[uint16(ext)] = pm[1]
			}
		}

		es := make(map[string]request.ManifestSpecTemplateContainerEnv)
		if len(env) > 0 {
			for _, e := range env {
				kv := strings.SplitN(e, "=", 2)
				eo := request.ManifestSpecTemplateContainerEnv{
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
				eo := request.ManifestSpecTemplateContainerEnv{
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
			senvs := make([]request.ManifestSpecTemplateContainerEnv, 0)
			for _, e := range es {
				senvs = append(senvs, e)
			}
			cs.Env = senvs
		}

		opts.Meta.Description = &description
		cs.Image.Name = image

		if secret != types.EmptyString {
			cs.Image.Secret = secret
		}

		css = append(css, cs)

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
