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
	"encoding/base64"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"io/ioutil"
	"os"
)

func init() {
	secretCreateCmd.Flags().StringP("text", "t", types.EmptyString, "write raw data")
	secretCreateCmd.Flags().StringArrayP("file", "f", make([]string, 0), "create secret from files")
	secretCreateCmd.Flags().BoolP("auth", "a", false, "create auth secret")
	secretCreateCmd.Flags().StringP("username", "u", types.EmptyString, "add username to registry secret")
	secretCreateCmd.Flags().StringP("password", "p", types.EmptyString, "add password to registry secret")
	secretCmd.AddCommand(secretCreateCmd)
}

const secretCreateExample = `
  # Create secret 'token' with 'secret' data 
  lb secret create token secret"
`

var secretCreateCmd = &cobra.Command{
	Use:     "create [NAME]",
	Short:   "Create secret",
	Example: secretCreateExample,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		auth, err := cmd.Flags().GetBool("auth")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		text, err := cmd.Flags().GetString("text")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		files, err := cmd.Flags().GetStringArray("file")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		name := args[0]
		opts := new(request.SecretCreateOptions)
		opts.Name = name
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

			username, err := cmd.Flags().GetString("username")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			password, err := cmd.Flags().GetString("password")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

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
			return
		}

		if err := opts.Validate(); err != nil {
			fmt.Println(err.Err())
			return
		}

		cli := envs.Get().GetClient()
		response, err := cli.V1().Secret().Create(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(fmt.Sprintf("Secret `%s` is created", name))

		secret := view.FromApiSecretView(response)
		secret.Print()
	},
}
