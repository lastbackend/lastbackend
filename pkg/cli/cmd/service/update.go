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
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

func UpdateCmd(name string, desc, sources *string, memory *int64) {

	ss, err := Update(name, desc, sources, memory)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fmt.Sprintf("Service `%s` is updated", name))

	ss.Print()
}

func Update(name string, desc, sources *string, memory *int64) (*view.Service, error) {

	stg := context.Get().GetStorage()
	cli := context.Get().GetClient()

	ns, err := stg.Namespace().Load()
	if err != nil {
		return nil, err
	}

	if ns == nil {
		return nil, errors.New("namespace has not been selected")
	}

	data := &request.ServiceUpdateOptions{
		Description: desc,
	}

	if sources == nil {
		data.Sources = sources
	}

	if memory == nil {
		data.Spec.Memory = memory
	}

	response, err := cli.V1().Namespace(ns.Meta.Name).Service(name).Update(context.Background(), data)
	if err != nil {
		return nil, err
	}

	ss := view.FromApiServiceView(response)

	return ss, nil
}
