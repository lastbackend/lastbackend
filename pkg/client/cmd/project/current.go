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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
)

func CurrentCmd() {

	var ctx = context.Get()

	project, err := Current()

	if err != nil {
		ctx.Log.Error(err)
		return
	}

	if project == nil {
		ctx.Log.Info("Project didn't select")
		return
	}

	project.DrawTable()
}

func Current() (*types.Project, error) {

	var (
		err     error
		ctx     = context.Get()
		project = new(types.Project)
	)

	if ctx.Token == "" {
		return nil, errors.NotLoggedMessage
	}

	err = ctx.Storage.Get("project", project)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return project, nil
}
