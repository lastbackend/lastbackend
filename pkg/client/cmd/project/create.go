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

type createS struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

func CreateCmd(name, description string) {

	var (
		log = c.Get().GetLogger()
	)

	err := Create(name, description)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Successful")
}

func Create(name, description string) error {

	var (
		err     error
		http    = c.Get().GetHttpClient()
		er      = new(errors.Http)
		project = new(types.Project)
	)

	if len(name) == 0 {
		return errors.BadParameter("name").Err()
	}

	_, _, err = http.
		POST("/project").
		AddHeader("Content-Type", "application/json").
		BodyJSON(createS{name, description}).
		Request(&project, er)
	if err != nil {
		return err
	}

	if er.Code == 401 {
		return errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	project, err = Switch(name)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
