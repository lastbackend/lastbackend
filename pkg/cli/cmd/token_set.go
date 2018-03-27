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

	"github.com/lastbackend/lastbackend/pkg/cli/storage"
	"github.com/spf13/cobra"
)

func init() {
	tokenCmd.AddCommand(tokenSetCmd)
}

const tokenSetExample = `
  # Set auth token for request quest in API 
  lb token set e3865d9b52c34dd4b6ec.5cff8c8e4cf6
`

var tokenSetCmd = &cobra.Command{
	Use:     "token [DATA]",
	Short:   "Set token to local storage",
	Example: tokenSetExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		token := args[0]

		if err := storage.SetToken(token); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Token successfully setted")
	},
}
