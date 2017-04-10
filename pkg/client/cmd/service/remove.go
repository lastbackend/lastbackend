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
	"github.com/lastbackend/lastbackend/pkg/errors"
)

func RemoveCmd(name string) {

	var (
		log = c.Get().GetLogger()
	)

	err := Remove(name)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Successful")
}

func Remove(name string) error {

	var (
		err     error
		http    = c.Get().GetHttpClient()
		service = new(types.Project)
		er      = new(errors.Http)
	)

	_, _, err = http.
		DELETE("/service/"+name).
		Request(service, er)

	if err != nil {
		return errors.New(err.Error())
	}

	if er.Code == 401 {
		return errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	return nil
}
