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
	"github.com/spf13/cobra"
)

func init() {
	serviceCmd.AddCommand(serviceRemoveCmd)
}

const serviceRemoveExample = `
  # Remove 'redis' service in 'ns-demo' namespace
  lb service remove ns-demo redis
`

var serviceRemoveCmd = &cobra.Command{
	Use:     "remove [NAMESPACE] [NAME]",
	Short:   "Remove service by name",
	Example: serviceRemoveExample,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]
		name := args[1]

		opts := &request.ServiceRemoveOptions{Force: false}

		if err := opts.Validate(); err != nil {
			fmt.Println(err.Err())
			return
		}

		cli := envs.Get().GetClient()
		cli.V1().Namespace(namespace).Service(name).Remove(envs.Background(), opts)

		fmt.Println(fmt.Sprintf("Service `%s` is successfully removed", name))
	},
}
