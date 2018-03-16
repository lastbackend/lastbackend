//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package app

import (
	"fmt"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
)

func RemoveCmd(name string) {

	if err := Remove(name); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fmt.Sprintf("App `%s` is successfully removed", name))
}

func Remove(name string) error {

	var (
		err     error
		http    = c.Get().GetHttpClient()
		storage = c.Get().GetStorage()
		er      = new(errors.Http)
		res     = new(struct{})
	)

	if len(name) == 0 {
		return errors.BadParameter("name").Err()
	}

	_, _, err = http.
		DELETE(fmt.Sprintf("/app/%s", name)).
		AddHeader("Content-Type", "application/json").
		Request(res, er)
	if err != nil {
		return errors.New(er.Message)
	}

	if er.Code == 401 {
		return errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	app, err := Current()
	if err != nil {
		return err
	}

	if app != nil && name == app.Meta.Name {
		if err := storage.App().Remove(); err != nil {
			return errors.UnknownMessage
		}
	}

	return nil
}
