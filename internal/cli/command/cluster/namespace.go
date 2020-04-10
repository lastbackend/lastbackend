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

	"github.com/lastbackend/lastbackend/internal/cli/views"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/util/decoder"
	v1 "github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/tools/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const namespaceListExample = `
  # Get all namespaces
  lb namespace ls"
`

const namespaceInspectExample = `
  # Get information for 'ns-demo' namespace
  lb namespace inspect ns-demo"
`

const namespaceCreateExample = `
  # Create 'ns-demo' namespace with description
  lb namespace create ns-demo --desc "Example description"
`

const namespaceRemoveExample = `
  # Remove 'ns-demo' namespace
  lb namespace remove ns-demo"
`

const namespaceUpdateExample = `
  # Update information for 'ns-demo' namespace
  lb namespace update ns-demo --desc "Example new description"
`

const namespaceApplyExample = `
  # Apply manifest from file or by URL
  lb namespace [name] apply -f"
`

func (c *command) NewNamespaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "namespace",
		Short: "Manage your namespaces",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			command := "[string]"
			if len(args) != 0 {
				command = args[0]
			}

			var ns = &cobra.Command{
				Use:   command,
				Short: "Manage your a namespace",
			}

			cmd.AddCommand(ns)

			if len(args) == 0 {
				if err := cmd.Help(); err != nil {
					log.Error(err.Error())
					return
				}
				return
			}

			// Attach sub command for namespace
			ns.AddCommand(
				c.NewServiceCmd(),
				c.NewSecretCmd(),
				c.NewRouteCmd(),
			)

			if err := ns.Execute(); err != nil {
				log.Error(err.Error())
				return
			}

		},
	}

	cmd.AddCommand(c.namespaceListCmd())
	cmd.AddCommand(c.namespaceInspectCmd())
	cmd.AddCommand(c.namespaceCreateCmd())
	cmd.AddCommand(c.namespaceRemoveCmd())
	cmd.AddCommand(c.namespaceUpdateCmd())
	cmd.AddCommand(c.namespaceApplyCmd())

	return cmd
}

func (c *command) namespaceListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "Display the namespace list",
		Example: namespaceListExample,
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			response, err := c.client.cluster.V1().Namespace().List(context.Background())

			if err != nil {
				fmt.Println(err)
				return
			}

			if response == nil || len(*response) == 0 {
				fmt.Println("no namespaces available")
				return
			}

			list := views.FromApiNamespaceListView(response)
			list.Print()
		},
	}

	return cmd
}

func (c *command) namespaceInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inspect [NAME]",
		Short:   "Get namespace info by name",
		Example: namespaceInspectExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]

			response, err := c.client.cluster.V1().Namespace(namespace).Get(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			ns := views.FromApiNamespaceView(response)
			ns.Print()
		},
	}

	return cmd
}

func (c *command) namespaceCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [NAME]",
		Short:   "Create new namespace",
		Example: namespaceCreateExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			opts := new(request.NamespaceManifest)

			desc := cmd.Flag("desc").Value.String()
			opts.Meta.Name = &args[0]
			opts.Meta.Description = &desc

			if err := opts.Validate(); err != nil {
				fmt.Println(err.Err())
				return
			}

			response, err := c.client.cluster.V1().Namespace().Create(context.Background(), opts)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("Namespace `%s` is created", opts.Meta.Name))
			ns := views.FromApiNamespaceView(response)
			ns.Print()
		},
	}

	cmd.Flags().StringP("desc", "d", "", "set namespace description (maximum 512 chars)")

	return cmd
}

func (c *command) namespaceRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove [NAME]",
		Short:   "Remove namespace by name",
		Example: namespaceRemoveExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]
			opts := &request.NamespaceRemoveOptions{Force: false}

			if err := opts.Validate(); err != nil {
				fmt.Println(err.Err())
				return
			}

			c.client.cluster.V1().Namespace(namespace).Remove(context.Background(), opts)

			fmt.Println(fmt.Sprintf("Namespace `%s` is successfully removed", namespace))
		},
	}

	return cmd
}

func (c *command) namespaceUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update [NAME]",
		Short:   "Update the namespace by name",
		Example: namespaceUpdateExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]
			desc := cmd.Flag("desc").Value.String()

			opts := new(request.NamespaceManifest)
			opts.Meta.Name = &namespace
			opts.Meta.Description = &desc

			if err := opts.Validate(); err != nil {
				fmt.Println(err.Err())
				return
			}

			response, err := c.client.cluster.V1().Namespace(namespace).Update(context.Background(), opts)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("Namespace `%s` is updated", namespace))
			ns := views.FromApiNamespaceView(response)
			ns.Print()
		},
	}

	cmd.Flags().StringP("desc", "d", "", "set namespace description (maximum 512 chars)")

	return cmd
}

func (c *command) namespaceApplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apply [NAME]",
		Short:   "Apply manifest files to cluster",
		Example: namespaceApplyExample,
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

			spec := v1.Request().Namespace().ApplyManifest()

			for _, f := range files {

				s, err := os.Open(f)
				if err != nil {
					if os.IsNotExist(err) {
						_ = fmt.Errorf("failed read data: file not exists: %s", f)
						os.Exit(1)
					}
				}

				if err := s.Close(); err != nil {
					fmt.Errorf("close file err: %s", err.Error())
					return
				}

				file, err := ioutil.ReadFile(f)
				if err != nil {
					_ = fmt.Errorf("failed read data from file: %s", f)
					os.Exit(1)
				}

				items := decoder.YamlSplit(file)
				fmt.Println("manifests:", len(items))

				for _, i := range items {

					var m = new(request.Runtime)

					if err := yaml.Unmarshal(i, m); err != nil {
						_ = fmt.Errorf("can not parse manifest: %s: %s", f, err.Error())
						continue
					}

					switch strings.ToLower(m.Kind) {
					case models.KindConfig:
						m := new(request.ConfigManifest)
						err := m.FromYaml(i)
						if err != nil {
							_ = fmt.Errorf("invalid specification: %s", err.Error())
							return
						}
						if m.Meta.Name == nil {
							break
						}
						fmt.Printf("Add config manifest: %s\n", *m.Meta.Name)
						spec.Configs[*m.Meta.Name] = m
						break
					case models.KindSecret:
						m := new(request.SecretManifest)
						err := m.FromYaml(i)
						if err != nil {
							_ = fmt.Errorf("invalid specification: %s", err.Error())
							return
						}
						if m.Meta.Name == nil {
							break
						}
						fmt.Printf("Add secret manifest: %s\n", *m.Meta.Name)
						spec.Secrets[*m.Meta.Name] = m
						break
					case models.KindService:
						m := new(request.ServiceManifest)
						err := m.FromYaml(i)
						if err != nil {
							_ = fmt.Errorf("invalid specification: %s", err.Error())
							return
						}
						if m.Meta.Name == nil {
							break
						}
						fmt.Printf("Add service manifest: %s\n", *m.Meta.Name)
						spec.Services[*m.Meta.Name] = m
						break
					case models.KindVolume:

						m := new(request.VolumeManifest)
						err := m.FromYaml(i)
						if err != nil {
							_ = fmt.Errorf("invalid specification: %s", err.Error())
							return
						}
						if m.Meta.Name == nil {
							break
						}
						fmt.Printf("Add volume manifest: %s\n", *m.Meta.Name)
						spec.Volumes[*m.Meta.Name] = m
						break
					case models.KindJob:
						m := new(request.JobManifest)
						err := m.FromYaml(i)
						if err != nil {
							_ = fmt.Errorf("invalid specification: %s", err.Error())
							return
						}
						if m.Meta.Name == nil {
							break
						}
						fmt.Printf("Add job manifest: %s\n", *m.Meta.Name)
						spec.Jobs[*m.Meta.Name] = m
						break
					case models.KindRoute:
						m := new(request.RouteManifest)
						err := m.FromYaml(i)
						if err != nil {
							_ = fmt.Errorf("invalid specification: %s", err.Error())
							return
						}
						if m.Meta.Name == nil {
							break
						}
						fmt.Printf("Add route manifest: %s\n", *m.Meta.Name)
						spec.Routes[*m.Meta.Name] = m
						break
					}
				}

				status, err := c.client.cluster.V1().Namespace(namespace).Apply(context.Background(), spec)
				if err != nil {
					_ = fmt.Errorf("invalid specification: %s", err.Error())
					return
				}

				fmt.Println()
				ns := views.FromApiNamespaceStatusView(status)
				ns.Print()
				return

			}
		},
	}

	cmd.Flags().StringArrayP("file", "f", make([]string, 0), "apply resources to namespace from files")

	return cmd
}
