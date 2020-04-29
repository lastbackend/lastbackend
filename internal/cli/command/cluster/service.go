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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lastbackend/lastbackend/tools/logger"
	"io"
	"os"
	"strings"

	"github.com/lastbackend/lastbackend/internal/cli/views"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/spf13/cobra"
)

const serviceListExample = `
  # Get all services for 'ns-demo' namespace  
  lb service ls ns-demo
`

const serviceInspectExample = `
  # Get information for 'redis' service in 'ns-demo' namespace
  lb service inspect ns-demo redis
`

const serviceCreateExample = `
  # Create new redis service with description and 256 MB limit memory
  lb service create ns-demo redis --desc "Example description" -m 256mib
`

const serviceRemoveExample = `
  # Remove 'redis' service in 'ns-demo' namespace
  lb service remove ns-demo redis
`

const serviceUpdateExample = `
  # Update info for 'redis' service in 'ns-demo' namespace
  lb service update ns-demo redis --desc "Example new description" -m 128
`

const serviceLogsExample = `
  # Get 'redis' service logs for 'ns-demo' namespace
  lb service logs ns-demo redis
`

func (c *command) NewServiceCmd() *cobra.Command {
	log := logger.WithContext(context.Background())
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage your service",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Error(err.Error())
				return
			}
		},
	}

	cmd.AddCommand(c.serviceListCmd())
	cmd.AddCommand(c.serviceInspectCmd())
	cmd.AddCommand(c.serviceCreateCmd())
	cmd.AddCommand(c.serviceRemoveCmd())
	cmd.AddCommand(c.serviceUpdateCmd())
	cmd.AddCommand(c.serviceLogsCmd())

	return cmd
}

func (c *command) serviceListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls [NAMESPACE]",
		Short:   "Display the services list",
		Example: serviceListExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]

			response, err := c.client.cluster.V1().Namespace(namespace).Service().List(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			if response == nil || len(*response) == 0 {
				fmt.Println("no services available")
				return
			}

			list := views.FromApiServiceListView(response)
			list.Print()
		},
	}
}

func (c *command) serviceInspectCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "inspect [NAMESPACE]/[NAME]",
		Short:   "Service info by name",
		Example: serviceInspectExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			namespace, name, err := serviceParseSelfLink(args[0])
			checkError(err)

			svc, err := c.client.cluster.V1().Namespace(namespace).Service(name).Get(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			routes, err := c.client.cluster.V1().Namespace(namespace).Route().List(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, r := range *routes {
				for _, rule := range r.Spec.Rules {
					if rule.Service == svc.Meta.Name {
						fmt.Println("exposed:", r.Status.State, r.Spec.Domain, r.Spec.Port)
					}
				}

			}

			ss := views.FromApiServiceView(svc)
			ss.Print()
		},
	}
}

func (c *command) serviceCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "create [NAMESPACE]/[NAME] [IMAGE]",
		Short:   "Create service",
		Example: serviceCreateExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			namespace, name, err := serviceParseSelfLink(args[0])
			checkError(err)

			image := args[1]

			opts, err := serviceParseManifest(cmd, name, image)
			checkError(err)

			response, err := c.client.cluster.V1().Namespace(namespace).Service().Create(context.Background(), opts)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("Service `%s` is created", name))

			service := views.FromApiServiceView(response)
			service.Print()
		},
	}
}

func (c *command) serviceRemoveCmd() *cobra.Command {
	return &cobra.Command{
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

			c.client.cluster.V1().Namespace(namespace).Service(name).Remove(context.Background(), opts)

			fmt.Println(fmt.Sprintf("Service `%s` is successfully removed", name))
		},
	}
}

func (c *command) serviceUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "update [NAMESPACE]/[NAME]",
		Short:   "Change configuration of the service",
		Example: serviceUpdateExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			namespace, name, err := serviceParseSelfLink(args[0])
			checkError(err)

			opts, err := serviceParseManifest(cmd, name, models.EmptyString)
			checkError(err)

			response, err := c.client.cluster.V1().Namespace(namespace).Service(name).Update(context.Background(), opts)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(fmt.Sprintf("Service `%s` is updated", name))
			ss := views.FromApiServiceView(response)
			ss.Print()
		},
	}
}

func (c *command) serviceLogsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "logs [NAMESPACE]/[NAME]",
		Short:   "Get service logs",
		Example: serviceLogsExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			opts := new(request.ServiceLogsOptions)

			var err error

			opts.Tail, err = cmd.Flags().GetInt("tail")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			opts.Follow, err = cmd.Flags().GetBool("follow")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			namespace, name, err := serviceParseSelfLink(args[0])
			checkError(err)

			reader, _, err := c.client.cluster.V1().Namespace(namespace).Service(name).Logs(context.Background(), opts)
			if err != nil {
				fmt.Println(err)
				return
			}

			dec := json.NewDecoder(reader)
			for {
				var doc models.LogMessage

				err := dec.Decode(&doc)
				if err == io.EOF {
					// all done
					break
				}
				if err != nil {
					fmt.Errorf(err.Error())
					os.Exit(1)
				}

				fmt.Println(">", doc.Selflink, doc.Data)
			}
		},
	}
}

func serviceParseSelfLink(selflink string) (string, string, error) {
	match := strings.Split(selflink, "/")

	var (
		namespace, name string
	)

	switch len(match) {
	case 2:
		namespace = match[0]
		name = match[1]
	case 1:
		fmt.Println("Use default namespace:", models.DEFAULT_NAMESPACE)
		namespace = models.DEFAULT_NAMESPACE
		name = match[0]
	default:
		return "", "", errors.New("invalid service name provided")
	}

	return namespace, name, nil
}

func serviceManifestFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("name", "n", "", "set service name")
	cmd.Flags().StringP("desc", "d", "", "set service description")
	cmd.Flags().StringP("memory", "m", "128MIB", "set service spec memory")
	cmd.Flags().IntP("replicas", "r", 0, "set service replicas")
	cmd.Flags().StringArrayP("port", "p", make([]string, 0), "set service ports")
	cmd.Flags().StringArrayP("env", "e", make([]string, 0), "set service env")
	cmd.Flags().StringArray("env-from-secret", make([]string, 0), "set service env from secret")
	cmd.Flags().StringArray("env-from-config", make([]string, 0), "set service env from config")
	cmd.Flags().StringP("image", "i", "", "set service image")
	cmd.Flags().String("image-secret-name", "", "set service image auth secret name")
	cmd.Flags().String("image-secret-key", "", "set service image auth secret key")
}

func serviceParseManifest(cmd *cobra.Command, name, image string) (*request.ServiceManifest, error) {

	var err error

	description, err := cmd.Flags().GetString("desc")
	checkFlagParseError(err)

	memory, err := cmd.Flags().GetString("memory")
	checkFlagParseError(err)

	if name == models.EmptyString {
		name, err = cmd.Flags().GetString("name")
		checkFlagParseError(err)
	}

	if image == models.EmptyString {
		image, err = cmd.Flags().GetString("image")
		checkFlagParseError(err)
	}

	ports, err := cmd.Flags().GetStringArray("ports")
	checkFlagParseError(err)

	env, err := cmd.Flags().GetStringArray("env")
	checkFlagParseError(err)

	senv, err := cmd.Flags().GetStringArray("env-from-secret")
	checkFlagParseError(err)

	cenv, err := cmd.Flags().GetStringArray("env-from-config")
	checkFlagParseError(err)

	replicas, err := cmd.Flags().GetInt("replicas")
	checkFlagParseError(err)

	authName, err := cmd.Flags().GetString("image-secret-name")
	checkFlagParseError(err)

	authKey, err := cmd.Flags().GetString("image-secret-key")
	checkFlagParseError(err)

	opts := new(request.ServiceManifest)
	css := make([]request.ManifestSpecTemplateContainer, 0)

	cs := request.ManifestSpecTemplateContainer{}

	if len(name) != 0 {
		opts.Meta.Name = &name
	}

	if len(description) != 0 {
		opts.Meta.Description = &description
	}

	if memory != models.EmptyString {
		cs.Resources.Request.RAM = memory
	}

	if replicas != 0 {
		opts.Spec.Replicas = &replicas
	}

	if len(ports) > 0 {
		opts.Spec.Network = new(request.ManifestSpecNetwork)
		opts.Spec.Network.Ports = make([]string, 0)
		opts.Spec.Network.Ports = ports
	}

	es := make(map[string]request.ManifestSpecTemplateContainerEnv)
	if len(env) > 0 {
		for _, e := range env {
			kv := strings.SplitN(e, "=", 2)
			eo := request.ManifestSpecTemplateContainerEnv{
				Name: kv[0],
			}
			if len(kv) > 1 {
				eo.Value = kv[1]
			}

			es[eo.Name] = eo
		}

	}
	if len(senv) > 0 {
		for _, e := range senv {
			kv := strings.SplitN(e, "=", 3)
			eo := request.ManifestSpecTemplateContainerEnv{
				Name: kv[0],
			}
			if len(kv) < 3 {
				return nil, errors.New("Service env from secret is in wrong format, should be [NAME]=[SECRET NAME]=[SECRET STORAGE KEY]")
			}

			if len(kv) == 3 {
				eo.Secret.Name = kv[1]
				eo.Secret.Key = kv[2]
			}

			es[eo.Name] = eo
		}
	}
	if len(cenv) > 0 {
		for _, e := range cenv {
			kv := strings.SplitN(e, "=", 3)
			eo := request.ManifestSpecTemplateContainerEnv{
				Name: kv[0],
			}
			if len(kv) < 3 {
				return nil, errors.New("Service env from config is in wrong format, should be [NAME]=[CONFIG NAME]=[CONFIG KEY]")
			}

			if len(kv) == 3 {
				eo.Config.Name = kv[1]
				eo.Config.Key = kv[2]
			}

			es[eo.Name] = eo
		}
	}

	if len(es) > 0 {
		senvs := make([]request.ManifestSpecTemplateContainerEnv, 0)
		for _, e := range es {
			senvs = append(senvs, e)
		}
		cs.Env = senvs
	}

	opts.Meta.Description = &description
	cs.Image.Name = image

	if authName != models.EmptyString {
		cs.Image.Secret.Name = authName
	}

	if authKey != models.EmptyString {
		cs.Image.Secret.Key = authKey
	}

	css = append(css, cs)

	if err := opts.Validate(); err != nil {
		return nil, err.Err()
	}

	return opts, nil
}
