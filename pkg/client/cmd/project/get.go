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

package project

import (
	p "github.com/lastbackend/lastbackend/pkg/client/api/views/v1/project"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
)

func GetCmd(name string) {

	var ctx = context.Get()

	project, err := Get(name)

	if err != nil {
		ctx.Log.Error(err)
		return
	}

	project.DrawTable()
}

func Get(name string) (*p.Project, error) {

	var (
		err     error
		ctx     = context.Get()
		er      = new(errors.Http)
		project = new(p.Project)
	)

	if len(name) == 0 {
		return nil, errors.BadParameter("name").Err()
	}

	_, _, err = ctx.HTTP.
		GET("/project/"+name).
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(&project, er)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if er.Code == 401 {
		return nil, errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, errors.New(er.Message)
	}

	return project, nil
}
