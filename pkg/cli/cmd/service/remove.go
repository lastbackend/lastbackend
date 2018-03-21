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
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

func RemoveCmd(name string) {

	if err := Remove(name); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fmt.Sprintf("Service `%s` is successfully removed", name))
}

func Remove(name string) error {

	stg := context.Get().GetStorage()
	cli := context.Get().GetClient()

	data := &request.ServiceRemoveOptions{
		Force: false,
	}

	ns, err := stg.Namespace().Load()
	if err != nil {
		return err
	}

	if ns == nil {
		return errors.New("namespace has not been selected")
	}

	return cli.V1().Namespace(ns.Meta.Name).Service(name).Remove(context.Background(), data)
}
