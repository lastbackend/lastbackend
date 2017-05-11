//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
	n "github.com/lastbackend/lastbackend/pkg/api/namespace/views/v1"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	nspace "github.com/lastbackend/lastbackend/pkg/cli/cmd/namespace"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
)

func UpdateCmd(name, nname, desc string, replicas int) {

	err := Update(name, nname, desc, replicas)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Service `" + name + "` successfully updated")
}

func Update(name, nname, desc string, replicas int) error {

	var (
		err  error
		http = c.Get().GetHttpClient()
		er   = new(errors.Http)
		ns   *n.Namespace
		res  = new(types.Namespace)
	)

	srv, _, err := Inspect(name)
	if err != nil {
		return errors.New(err.Error())
	}

	if nname == "" {
		nname = srv.Meta.Name
	}

	if desc == "" {
		desc = srv.Meta.Description
	}

	if replicas == 0 {
		replicas = srv.Meta.Replicas
	}

	cfg := types.ServiceUpdateConfig{
		Name:        &nname,
		Description: &desc,
		Replicas:    &replicas,
	}

	ns, err = nspace.Current()
	if err != nil {
		return errors.New(err.Error())
	}
	if ns == nil {
		return errors.New("Namespace didn't select")
	}

	_, _, err = http.
		PUT(fmt.Sprintf("/namespace/%s/service/%s", ns.Meta.Name, name)).
		AddHeader("Content-Type", "application/json").
		BodyJSON(cfg).
		Request(&res, er)
	if err != nil {
		return errors.New(err.Error())
	}

	if er.Code == 401 {
		return errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	return nil
}
