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

	"github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/cli/view"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

func ListCmd() {

	list, err := List()
	if err != nil {
		fmt.Println(err)
		return
	}

	list.Print()
}

func List() (*view.ServiceList, error) {

	stg := context.Get().GetStorage()
	cli := context.Get().GetClient()

	ns, err := stg.Namespace().Load()
	if err != nil {
		return nil, err
	}

	if ns == nil {
		return nil, errors.New("namespace has not been selected")
	}

	response, err := cli.V1().Namespace(ns.Meta.Name).Service().List(context.Background())
	if err != nil {
		return nil, err
	}

	ss := view.FromApiServiceListView(response)

	return ss, nil
}
