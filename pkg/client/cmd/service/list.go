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
	s "github.com/lastbackend/lastbackend/pkg/apis/views/v1/service"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	c "github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
)

func ListServiceCmd() {

	var (
		log = c.Get().GetLogger()
	)

	services, projectName, err := List()
	if err != nil {
		log.Error(err)
		return
	}

	if services != nil {
		services.DrawTable(projectName)
	}
}

func List() (*s.ServiceList, string, error) {

	var (
		err      error
		log      = c.Get().GetLogger()
		http     = c.Get().GetHttpClient()
		er       = new(errors.Http)
		services = new(s.ServiceList)
	)

	p, err := project.Current()
	if err != nil {
		return nil, "", errors.New(err.Error())
	}

	if p == nil {
		log.Info("Project didn't select")
		return nil, "", nil
	}

	_, _, err = http.
		GET("/p/"+p.Name+"/service").
		Request(services, er)
	if err != nil {
		return nil, "", errors.New(err.Error())
	}

	if er.Code == 401 {
		return nil, "", errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, "", errors.New(er.Message)
	}

	if len(*services) == 0 {
		log.Info("You don't have any services")
		return nil, "", nil
	}

	return services, p.Name, nil
}
