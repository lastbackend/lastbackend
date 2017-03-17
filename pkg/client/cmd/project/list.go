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
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func ListProjectCmd() {

	var ctx = context.Get()

	projects, err := List()
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	if projects != nil {
		projects.DrawTable()
	}
}

func List() (*model.ProjectList, error) {

	var (
		err      error
		ctx      = context.Get()
		er       = new(e.Http)
		projects = new(model.ProjectList)
	)

	_, _, err = ctx.HTTP.
		GET("/project").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(projects, er)
	if err != nil {
		return nil, err
	}

	if er.Code == 401 {
		return nil, e.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, errors.New(er.Message)
	}

	if len(*projects) == 0 {
		ctx.Log.Info("You don't have any projects")
		return nil, nil
	}

	return projects, nil
}
