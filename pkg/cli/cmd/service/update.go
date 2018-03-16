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
	av "github.com/lastbackend/lastbackend/pkg/api/app/views/v1"
	a "github.com/lastbackend/lastbackend/pkg/cli/cmd/app"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"github.com/lastbackend/lastbackend/pkg/common/types"
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
		app  *av.App
		res  = new(types.App)
	)

	srv, _, err := Inspect(name)
	if err != nil {
		return err
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

	app, err = a.Current()
	if err != nil {
		return err
	}
	if app == nil {
		return errors.New("App didn't select")
	}

	_, _, err = http.
		PUT(fmt.Sprintf("/app/%s/service/%s", app.Meta.Name, name)).
		AddHeader("Content-Type", "application/json").
		BodyJSON(cfg).
		Request(&res, er)
	if err != nil {
		return err
	}

	if er.Code == 401 {
		return errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	return nil
}
