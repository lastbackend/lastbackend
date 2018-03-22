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

package namespace

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/cli/envs"
	"github.com/lastbackend/lastbackend/pkg/cli/view"
)

func UpdateCmd(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		cmd.Help()
		return
	}

	name := args[0]
	desc := cmd.Flag("desc").Value.String()

	opts := new(request.NamespaceUpdateOptions)
	opts.Description = &desc

	if err := opts.Validate(); err != nil {
		fmt.Println(err.Attr)
		return
	}

	cli := envs.Get().GetClient()
	response, err := cli.V1().Namespace(name).Update(envs.Background(), opts)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fmt.Sprintf("Namespace `%s` is updated", name))
	ns := view.FromApiNamespaceView(response)
	ns.Print()
}
