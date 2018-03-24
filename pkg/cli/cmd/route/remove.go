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

package route

import (
	"fmt"

	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/cli/envs"
	"github.com/spf13/cobra"
)

func RemoveCmd(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		cmd.Help()
		return
	}
	name := args[0]

	namespace, _ := cmd.Flags().GetString("namespace")

	if namespace == "" {
		fmt.Println("namesapace parameter not set")
		return
	}

	opts := &request.RouteRemoveOptions{Force: false}

	if err := opts.Validate(); err != nil {
		fmt.Println(err.Attr)
		return
	}

	cli := envs.Get().GetClient()
	cli.V1().Namespace(namespace).Route(name).Remove(envs.Background(), opts)

	fmt.Println(fmt.Sprintf("Route `%s` remove now", name))
}
