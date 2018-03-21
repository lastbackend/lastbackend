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
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package service

import (
	"fmt"

	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/cli/envs"
	"github.com/lastbackend/lastbackend/pkg/cli/view"
	"github.com/spf13/cobra"
)

func CreateCmd(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		cmd.Help()
		return
	}

	namespace, _ := cmd.Flags().GetString("namespace")

	if namespace == "" {
		fmt.Println("namesapace parameter not set")
		return
	}

	name := args[0]

	description, _ := cmd.Flags().GetString("desc")
	memory, _ := cmd.Flags().GetInt64("memory")
	image, _ := cmd.Flags().GetString("image")
	replicas, _ := cmd.Flags().GetInt("replicas")

	opts := new(request.ServiceCreateOptions)
	opts.Spec = new(request.ServiceOptionsSpec)

	opts.Name = &name
	opts.Description = &description
	opts.Spec.Memory = &memory
	opts.Image = &image
	opts.Replicas = &replicas

	if err := opts.Validate(); err != nil {
		fmt.Println(err.Attr)
		return
	}

	cli := envs.Get().GetClient()
	response, err := cli.V1().Namespace(namespace).Service().Create(envs.Background(), opts)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fmt.Sprintf("Service `%s` is created", name))

	service := view.FromApiServiceView(response)
	service.Print()
}
