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
	"github.com/lastbackend/lastbackend/libs/editor"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"gopkg.in/yaml.v2"
	"strings"
)

func UpdateCmd(name string) {

	var ctx = context.Get()

	serviceModel, _, err := Inspect(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	config, err := GetConfig(serviceModel)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	err = Update(name, *config)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	ctx.Log.Info("Successful")
}

func Update(name string, config model.ServiceUpdateConfig) error {

	var (
		err     error
		ctx     = context.Get()
		er      = new(e.Http)
		project = new(model.Project)
		res     = new(model.Project)
	)

	err = ctx.Storage.Get("project", project)
	if err != nil {
		return errors.New(err.Error())
	}

	if project.ID == "" {
		return errors.New("Project didn't select")
	}

	_, _, err = ctx.HTTP.
		PUT("/project/"+project.ID+"/service/"+name).
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		BodyJSON(config).
		Request(&res, er)

	if err != nil {
		return err
	}

	if er.Code == 401 {
		return e.NotLoggedMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	return nil
}

func GetConfig(service *model.Service) (*model.ServiceUpdateConfig, error) {

	var config = service.GetConfig()

	buf, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	// TODO: To allow for the possibility of naming the session re-editing
	res, err := editor.Run(strings.NewReader(string(buf)))
	if err != nil {
		return nil, err
	}

	err = res.FromYAML(&config)
	if err != nil {
		// TODO: When is have error parse yaml. Ask question about reopen config for correct this
		return nil, err
	}

	return config, nil
}
