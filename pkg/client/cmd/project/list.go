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

func ListProjectCmd() {

	var (
		log = c.Get().GetLogger()
	)

	projects, err := List()
	if err != nil {
		log.Error(err)
		return
	}

	if projects != nil {
		projects.DrawTable()
	}
}

func List() (*p.ProjectList, error) {

	var (
		err      error
		log      = c.Get().GetLogger()
		http     = c.Get().GetHttpClient()
		er       = new(errors.Http)
		projects = new(p.ProjectList)
	)

	_, _, err = http.
		GET("/project").
		Request(projects, er)
	if err != nil {
		return nil, err
	}

	if er.Code == 401 {
		return nil, errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, errors.New(er.Message)
	}

	if len(*projects) == 0 {
		log.Info("You don't have any projects")
		return nil, nil
	}

	return projects, nil
}
