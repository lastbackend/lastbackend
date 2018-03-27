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
	"github.com/lastbackend/lastbackend/pkg/cli/view"
	"github.com/spf13/cobra"
)

func init() {
	secretCmd.AddCommand(secretListCmd)
}

const secretListExample = `
  # Get all secrets records
  lb secret ls ns-demo"
`

var secretListCmd = &cobra.Command{
	Use:     "ls [NAMESPACE]",
	Short:   "Display the secrets list",
	Example: secretListExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]

		cli := envs.Get().GetClient()
		response, err := cli.V1().Namespace(namespace).Secret().List(envs.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		if response == nil || len(*response) == 0 {
			fmt.Println("no secrets available")
			return
		}

		list := view.FromApiSecretListView(response)
		list.Print()
	},
}
