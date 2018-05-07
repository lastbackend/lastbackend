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
	"github.com/lastbackend/lastbackend/pkg/cli/envs"
	"github.com/spf13/cobra"
	"github.com/lastbackend/lastbackend/pkg/cli/view"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

func init() {
	serviceCmd.AddCommand(serviceWatchCmd)
}

const serviceWatchExample = `
  # Get 'redis' service watch for 'ns-demo' namespace  
  lb service watch ns-demo redis
`

var serviceWatchCmd = &cobra.Command{
	Use:     "watch [NAMESPACE] [NAME]",
	Short:   "Get service watch",
	Example: serviceWatchExample,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]
		name := args[1]

		cli := envs.Get().GetClient()

		watcher, err := cli.V1().Namespace(namespace).Service(name).Watch(envs.Background())
		if err != nil {
			fmt.Println(err)
			return
		}
		defer watcher.Stop()

		for w := range watcher.ResultChan() {
			fmt.Println(">>>>>>>>>", w.Data)
			if w.Data == nil {
				continue
			}

			service := w.Data.(*views.Service)
			ss := view.FromApiServiceView(service)
			ss.Print()
		}

	},
}
