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

	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/cli/view"
)

func CreateCmd(name, desc string) {

	ns, err := Create(name, desc)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fmt.Sprintf("Namespace `%s` is created", name))

	ns.Print()
}

func Create(name, desc string) (*view.Namespace, error) {

	stg := context.Get().GetStorage()
	cli := context.Get().GetClient()

	data := &request.NamespaceCreateOptions{
		Name:        name,
		Description: desc,
	}

	response, err := cli.V1().Namespace().Create(context.Background(), data)
	if err != nil {
		return nil, err
	}

	ns := view.FromApiNamespaceView(response)

	if err := stg.Namespace().Save(ns); err != nil {
		return nil, err
	}

	return ns, nil
}
