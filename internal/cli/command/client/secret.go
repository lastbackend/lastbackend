//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package client

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/tools/logger"
	"io/ioutil"
	"os"
	"strings"

	"github.com/lastbackend/lastbackend/internal/cli/views"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/spf13/cobra"
)

const secretListExample = `
  # Get all secrets records
  lb secret ls"
`

const secretInspectExample = `
  # Inspect secret 'token' 
  lb secret inspect token"
`

const secretCreateExample = `
  # Create secret 'token' with 'secret' data 
  lb secret create token secret"
`

const secretRemoveExample = `
  # Remove 'token' secret
  lb secret remove token
`

const secretUpdateExample = `
  # Update 'token' secret record with 'new-secret' data
  lb secret update token new-secret"
`

func (c *command) NewSecretCmd() *cobra.Command {
	log := logger.WithContext(context.Background())
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "Manage your secret",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Error(err.Error())
				return
			}
		},
	}

	cmd.AddCommand(c.secretListCmd())
	cmd.AddCommand(c.secretInspectCmd())
	cmd.AddCommand(c.secretCreateCmd())
	cmd.AddCommand(c.secretRemoveCmd())
	cmd.AddCommand(c.secretUpdateCmd())

	return cmd
}

func (c *command) secretListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls [NAMESPACE]",
		Short:   "Display the secrets list",
		Example: secretListExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]

			response, err := c.client.cluster.V1().Namespace(namespace).Secret().List(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			if response == nil || len(*response) == 0 {
				fmt.Println("no secrets available")
				return
			}

			list := views.FromApiSecretListView(response)
			list.Print()
		},
	}
}

func (c *command) secretInspectCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "inspect [NAMESPACE] [NAME]",
		Short:   "Inspect secret",
		Example: secretInspectExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]
			name := args[1]

			response, err := c.client.cluster.V1().Namespace(namespace).Secret(name).Get(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			secret := views.FromApiSecretView(response)
			secret.Print()
		},
	}
}

func (c *command) secretCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "create [NAMESPACE] [NAME]",
		Short:   "Create secret",
		Example: secretCreateExample,
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			auth, err := cmd.Flags().GetBool("auth")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			text, err := cmd.Flags().GetStringArray("text")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			files, err := cmd.Flags().GetStringArray("file")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			namespace := args[0]
			name := args[1]
			opts := new(request.SecretManifest)
			opts.Meta.Name = &name
			opts.Spec.Data = make(map[string]string, 0)

			switch true {
			case len(text) > 0:
				opts.Spec.Type = models.KindSecretOpaque

				for _, t := range text {
					var (
						k, v string
					)

					kv := strings.SplitN(t, "=", 2)
					k = kv[0]
					if len(kv) > 1 {
						v = kv[1]
					}
					opts.Spec.Data[k] = v
				}

				break
			case auth:
				opts.Spec.Type = models.KindSecretAuth
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
				opts.SetAuthData(username, password)
				break
			case len(files) > 0:
				opts.Spec.Type = models.KindSecretOpaque
				for _, f := range files {
					c, err := ioutil.ReadFile(f)
					if err != nil {
						_ = fmt.Errorf("failed read data from file: %s", f)
						os.Exit(1)
					}
					opts.Spec.Data[f] = string(c)
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

			response, err := c.client.cluster.V1().Namespace(namespace).Secret().Create(context.Background(), opts)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("Secret `%s` is created", name))

			secret := views.FromApiSecretView(response)
			secret.Print()
		},
	}
}

func (c *command) secretRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove [NAMESPACE] [NAME]",
		Short:   "Remove secret by name",
		Example: secretRemoveExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]
			name := args[1]

			opts := &request.SecretRemoveOptions{Force: false}

			if err := opts.Validate(); err != nil {
				fmt.Println(err.Err())
				return
			}

			c.client.cluster.V1().Namespace(namespace).Secret(name).Remove(context.Background(), opts)

			fmt.Println(fmt.Sprintf("Secret `%s` remove now", name))
		},
	}
}

func (c *command) secretUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "update [NAME]",
		Short:   "Change configuration of the secret",
		Example: secretUpdateExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			auth, _ := cmd.Flags().GetBool("auth")
			text, _ := cmd.Flags().GetStringArray("text")
			files, _ := cmd.Flags().GetStringArray("file")

			namespace := args[0]
			name := args[1]
			opts := new(request.SecretManifest)
			opts.Spec.Data = make(map[string]string, 0)

			switch true {
			case len(text) > 0:
				opts.Spec.Type = models.KindSecretOpaque

				for _, t := range text {
					var (
						k, v string
					)

					kv := strings.SplitN(t, "=", 2)
					k = kv[0]
					if len(kv) > 1 {
						v = kv[1]
					}
					opts.Spec.Data[k] = v
				}

				break
			case auth:
				opts.Spec.Type = models.KindSecretAuth

				username, _ := cmd.Flags().GetString("username")
				password, _ := cmd.Flags().GetString("password")

				opts.SetAuthData(username, password)

				break
			case len(files) > 0:
				opts.Spec.Type = models.KindSecretOpaque
				for _, f := range files {
					c, err := ioutil.ReadFile(f)
					if err != nil {
						_ = fmt.Errorf("failed read data from file: %s", f)
						os.Exit(1)
					}
					opts.Spec.Data[f] = string(c)
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

			response, err := c.client.cluster.V1().Namespace(namespace).Secret(name).Update(context.Background(), opts)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("Secret `%s` is updated", name))
			ss := views.FromApiSecretView(response)
			ss.Print()
		},
	}
}
