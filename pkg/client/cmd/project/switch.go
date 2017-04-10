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
	c "github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
)

func SwitchCmd(name string) {

	var (
		log = c.Get().GetLogger()
	)

	project, err := Switch(name)
	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("The project `%s` was selected as the current", project.Meta.Name)
}

func Switch(name string) (*types.Project, error) {

	var (
		er      = new(errors.Http)
		http    = c.Get().GetHttpClient()
		storage = c.Get().GetStorage()
		project = new(types.Project)
	)

	_, _, err := http.
		GET("/project/"+name).
		AddHeader("Content-Type", "application/json").
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

	err = storage.Set("project", project)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return project, nil
}
