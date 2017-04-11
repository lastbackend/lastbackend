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
	c "github.com/lastbackend/lastbackend/pkg/client/context"
	p "github.com/lastbackend/lastbackend/pkg/daemon/api/views/v1/project"
	"github.com/lastbackend/lastbackend/pkg/errors"
)

func CurrentCmd() {

	var (
		log = c.Get().GetLogger()
	)

	project, err := Current()

	if err != nil {
		log.Error(err)
		return
	}

	if project == nil {
		log.Info("Project didn't select")
		return
	}

	project.DrawTable()
}

func Current() (*p.Project, error) {

	var (
		err     error
		storage = c.Get().GetStorage()
		project = new(p.Project)
	)

	err = storage.Get("project", project)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return project, nil
}
