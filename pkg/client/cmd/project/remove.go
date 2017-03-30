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

func RemoveCmd(name string) {

	var ctx = context.Get()

	err := Remove(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	ctx.Log.Info("Successful")
}

func Remove(name string) error {

	var (
		err error
		ctx = context.Get()
		er  = new(errors.Http)
		res = new(struct{})
	)

	if len(name) == 0 {
		return errors.BadParameter("name").Err()
	}

	_, _, err = ctx.HTTP.
		DELETE("/project/"+name).
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(res, er)
	if err != nil {
		return errors.New(err.Error())
	}

	if er.Code == 401 {
		return errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	project, err := Current()
	if err != nil {
		return errors.New(err.Error())
	}

	if project != nil {
		if name == project.Name {
			err = ctx.Storage.Set("project", types.Project{})
			if err != nil {
				return errors.New(err.Error())
			}
		}
	}

	return nil
}
