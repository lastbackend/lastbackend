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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"encoding/base64"
	"io/ioutil"
	"os"
)

func init() {
	secretUpdateCmd.Flags().StringP("text", "t", "", "write raw data")
	secretUpdateCmd.Flags().StringArrayP("file", "f", make([]string, 0), "create secret from files")
	secretUpdateCmd.Flags().BoolP("auth", "a", false, "create auth secret")
	secretUpdateCmd.Flags().StringP("username", "u", types.EmptyString, "add username to registry secret")
	secretUpdateCmd.Flags().StringP("password", "p", types.EmptyString, "add password to registry secret")
	secretCmd.AddCommand(secretUpdateCmd)
}

const secretUpdateExample = `
  # Update 'token' secret record with 'new-secret' data
  lb secret update token new-secret"
`

var secretUpdateCmd = &cobra.Command{
	Use:     "update [NAME]",
	Short:   "Change configuration of the secret",
	Example: secretUpdateExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		auth, _ := cmd.Flags().GetBool("auth")
		text, _ := cmd.Flags().GetString("text")
		files, _ := cmd.Flags().GetStringArray("file")

		name := args[0]
		opts := new(request.SecretUpdateOptions)
		opts.Data = make(map[string][]byte, 0)

		switch true {
		case text != types.EmptyString:
			var (
				data = []byte(text)
			)
			opts.Kind = types.KindSecretText
			opts.Data[types.KindSecretText] = []byte(base64.StdEncoding.EncodeToString(data))
			break
		case auth:
			opts.Kind = types.KindSecretAuth

			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")

			s := new(types.Secret)
			s.EncodeSecretAuthData(types.SecretAuthData{
				Username: username,
				Password: password,
			})
			opts.Data = s.Data

			break
		case len(files) > 0:
			opts.Kind = types.KindSecretFiles
			for _, f := range files {
				c, err := ioutil.ReadFile(f)
				if err != nil {
					_ = fmt.Errorf("failed read data from file: %s", f)
					os.Exit(1)
				}
				opts.Data[f] = c
			}
			break
		default:
			fmt.Println("You need to provide secret type")
			os.Exit(0)
		}

		if err := opts.Validate(); err != nil {
			fmt.Println(err.Err())
			return
		}

		cli := envs.Get().GetClient()
		response, err := cli.V1().Secret(name).Update(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(fmt.Sprintf("Secret `%s` is updated", name))
		ss := view.FromApiSecretView(response)
		ss.Print()
	},
}
