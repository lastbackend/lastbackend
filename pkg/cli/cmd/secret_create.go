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
	secretCmd.AddCommand(secretCreateCmd)
}

const secretCreateExample = `
  # Create secret 'token' for 'ns-demo' namespace with 'secret' data 
  lb secret create ns-demo token secret"
`

var secretCreateCmd = &cobra.Command{
	Use:     "create [NAMESPACE] [NAME] [DATA]",
	Short:   "Create secret",
	Example: secretCreateExample,
	Args:    cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]
		name := args[1]
		data := args[2]

		opts := new(request.SecretCreateOptions)
		opts.Data = &data

		if err := opts.Validate(); err != nil {
			fmt.Println(err.Err())
			return
		}

		cli := envs.Get().GetClient()
		response, err := cli.V1().Namespace(namespace).Secret().Create(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(fmt.Sprintf("Secret `%s` is created", name))

		service := view.FromApiSecretView(response)
		service.Print()
	},
}
