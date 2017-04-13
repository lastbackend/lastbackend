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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	c "github.com/lastbackend/lastbackend/pkg/client/context"
	s "github.com/lastbackend/lastbackend/pkg/daemon/service/views/v1"
	"github.com/lastbackend/lastbackend/pkg/editor"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"gopkg.in/yaml.v2"
	"strings"
)

func UpdateCmd(name string) {

	var (
		log = c.Get().GetLogger()
	)

	serviceModel, _, err := Inspect(name)
	if err != nil {
		log.Error(err)
		return
	}

	config, err := GetConfig(serviceModel)
	if err != nil {
		log.Error(err)
		return
	}

	err = Update(name, *config)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Successful")
}

func Update(name string, config types.ServiceUpdateConfig) error {

	var (
		err       error
		http      = c.Get().GetHttpClient()
		storage   = c.Get().GetStorage()
		er        = new(errors.Http)
		namespace = new(types.Namespace)
		res       = new(types.Namespace)
	)

	err = storage.Get("namespace", namespace)
	if err != nil {
		return errors.New(err.Error())
	}

	if namespace.Meta.Name == "" {
		return errors.New("Namespace didn't select")
	}

	_, _, err = http.
		PUT("/namespace/"+namespace.Meta.Name+"/service/"+name).
		AddHeader("Content-Type", "application/json").
		BodyJSON(config).
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

func GetConfig(service *s.Service) (*types.ServiceUpdateConfig, error) {

	var config = struct{}{}

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

	return &types.ServiceUpdateConfig{}, nil
}
