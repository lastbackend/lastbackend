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
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	p "github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func ListServiceCmd() {

	var ctx = context.Get()

	services, projectName, err := List()
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	if services != nil {
		services.DrawTable(projectName)
	}
}

func List() (*model.ServiceList, string, error) {

	var (
		err      error
		ctx      = context.Get()
		er       = new(e.Http)
		services = new(model.ServiceList)
	)

	project, err := p.Current()
	if err != nil {
		return nil, "", errors.New(err.Error())
	}

	if project == nil {
		ctx.Log.Info("Project didn't select")
		return nil, "", nil
	}

	_, _, err = ctx.HTTP.
		GET("/project/"+project.Name+"/service").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(services, er)
	if err != nil {
		return nil, "", errors.New(err.Error())
	}

	if er.Code == 401 {
		return nil, "", e.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, "", errors.New(er.Message)
	}

	if len(*services) == 0 {
		ctx.Log.Info("You don't have any services")
		return nil, "", nil
	}

	return services, project.Name, nil
}
