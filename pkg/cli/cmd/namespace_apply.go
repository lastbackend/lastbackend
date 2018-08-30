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
	"io/ioutil"
	"os"

	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/cli/envs"
	"github.com/lastbackend/lastbackend/pkg/cli/view"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
)

const applyExample = `
  # Apply manifest from file or by URL
  lb apply -f"
`

func init() {
	applyCmd.Flags().StringArrayP("file", "f", make([]string, 0), "create secret from files")
	namespaceCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply [NAME]",
	Short: "Apply file manifest to cluster",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]

		files, err := cmd.Flags().GetStringArray("file")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if len(files) == 0 {
			cmd.Help()
		}

		cli := envs.Get().GetClient()

		for _, f := range files {
			c, err := ioutil.ReadFile(f)
			if err != nil {
				_ = fmt.Errorf("failed read data from file: %s", f)
				os.Exit(1)
			}

			var m = new(request.Runtime)
			yaml.Unmarshal(c, m)

			if m.Kind == "Service" {

				spec := v1.Request().Service().Manifest()
				err := spec.FromYaml(c)
				if err != nil {
					fmt.Errorf("invalid specification: %s", err.Error())
					return
				}

				var rsvc *views.Service

				if m.Meta.Name != nil {
					rsvc, _ = cli.V1().Namespace(namespace).Service(*m.Meta.Name).Get(envs.Background())
				}

				if rsvc == nil {
					fmt.Println("create new service")
					rsvc, err = cli.V1().Namespace(namespace).Service().Create(envs.Background(), spec)
					if err != nil {
						fmt.Println(err)
						return
					}
				} else {
					rsvc, err = cli.V1().Namespace(namespace).Service(rsvc.Meta.Name).Update(envs.Background(), spec)
					if err != nil {
						fmt.Println(3)
						fmt.Println(err)
						return
					}
				}

				if rsvc != nil {
					service := view.FromApiServiceView(rsvc)
					service.Print()
				} else {
					fmt.Println("ooops")
				}

			}

			return

		}
	},
}
