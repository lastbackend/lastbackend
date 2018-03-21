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
	"github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/cli/view"
	"github.com/spf13/cobra"
)

func CreateCmd(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		cmd.Help()
		return
	}
	var name = args[0]

	opts := &request.ServiceCreateOptions{Name: &name}
	cmd.Flags().StringVarP(opts.Description, "desc", "d", "", "Set description")
	cmd.Flags().StringVarP(opts.Sources, "sources", "s", "", "Set sources")
	cmd.Flags().Int64VarP(opts.Spec.Memory, "memory", "m", 0, "Set memory")

	cli := context.Get().GetClient()
	response, err := cli.V1().Namespace(name).Service().Create(context.Background(), opts)
	if err != nil {
		return
	}

	fmt.Println(fmt.Sprintf("Service `%s` is created", name))
	service := view.FromApiServiceView(response)
	service.Print()
}
