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

package cluster

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/lastbackend/lastbackend/tools/logger"
	"github.com/lastbackend/lastbackend/internal/cli/views"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/spf13/cobra"
)

const configListExample = `
  # Get all configs records
  lb config ls"
`

const configInspectExample = `
  # Inspect config 'name' 
  lb config inspect namespace name"
`

const configCreateExample = `
  # Create config 'token' with 'config' data 
  lb config create token config"
`

const configRemoveExample = `
  # Remove 'name' config in namespace
  lb config remove namespace name
`

const configUpdateExample = `
  # Update 'token' config record with 'new-config' data
  lb config update token new-config"
`

func (c *command) NewConfigCmd() *cobra.Command {
	log := logger.WithContext(context.Background())

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage your configs",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Error(err.Error())
				return
			}
		},
	}

	cmd.AddCommand(c.configListCmd())
	cmd.AddCommand(c.configInspectCmd())
	cmd.AddCommand(c.configCreateCmd())
	cmd.AddCommand(c.configRemoveCmd())
	cmd.AddCommand(c.configUpdateCmd())

	return cmd
}

func (c *command) configListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls [NAMESPACE]",
		Short:   "Display the configs list",
		Example: configListExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]

			response, err := c.client.cluster.V1().Namespace(namespace).Config().List(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			if response == nil || len(*response) == 0 {
				fmt.Println("no configs available")
				return
			}

			list := views.FromApiConfigListView(response)
			list.Print()
		},
	}

	return cmd
}

func (c *command) configInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inspect [NAMESPACE] [NAME]",
		Short:   "Inspect config",
		Example: configInspectExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]
			name := args[1]

			response, err := c.client.cluster.V1().Namespace(namespace).Config(name).Get(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			config := views.FromApiConfigView(response)
			config.Print()
		},
	}

	return cmd
}

func (c *command) configCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [NAME]",
		Short:   "Create config",
		Example: configCreateExample,
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

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

			opts := new(request.ConfigManifest)
			opts.Meta.Name = &name
			opts.Spec.Data = make(map[string]string, 0)

			switch true {
			case len(text) > 0:
				opts.Spec.Type = models.KindConfigText

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
			case len(files) > 0:
				opts.Spec.Type = models.KindConfigText
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
				fmt.Println("You need to provide config type")
				return
			}

			if err := opts.Validate(); err != nil {
				fmt.Println(err.Err())
				return
			}

			response, err := c.client.cluster.V1().Namespace(namespace).Config().Create(context.Background(), opts)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("Secret `%s` is created", name))

			config := views.FromApiConfigView(response)
			config.Print()
		},
	}

	cmd.Flags().StringArrayP("text", "t", make([]string, 0), "write text data in key=value format")
	cmd.Flags().StringArrayP("file", "f", make([]string, 0), "create config from files")

	return cmd
}

func (c *command) configRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove [NAMESPACE] [NAME]",
		Short:   "Remove config by name",
		Example: configRemoveExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]
			name := args[1]

			opts := &request.ConfigRemoveOptions{Force: false}

			if err := opts.Validate(); err != nil {
				fmt.Println(err.Err())
				return
			}

			c.client.cluster.V1().Namespace(namespace).Config(name).Remove(context.Background(), opts)

			fmt.Println(fmt.Sprintf("Config `%s` remove now", name))
		},
	}

	return cmd
}

func (c *command) configUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update [NAME]",
		Short:   "Change configuration of the config",
		Example: configUpdateExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			text, _ := cmd.Flags().GetStringArray("text")
			files, _ := cmd.Flags().GetStringArray("file")

			namespace := args[0]
			name := args[1]
			opts := new(request.ConfigManifest)
			opts.Spec.Data = make(map[string]string, 0)

			switch true {
			case len(text) > 0:
				opts.Spec.Type = models.KindConfigText

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
			case len(files) > 0:
				opts.Spec.Type = models.KindConfigText
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
				fmt.Println("You need to provide config type")
				return
			}

			if err := opts.Validate(); err != nil {
				fmt.Println(err.Err())
				return
			}

			response, err := c.client.cluster.V1().Namespace(namespace).Config(name).Update(context.Background(), opts)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("Config `%s` is updated", name))
			ss := views.FromApiConfigView(response)
			ss.Print()
		},
	}

	cmd.Flags().StringArrayP("text", "t", make([]string, 0), "write config in key=value format")
	cmd.Flags().StringArrayP("file", "f", make([]string, 0), "create config from files")

	return cmd
}
