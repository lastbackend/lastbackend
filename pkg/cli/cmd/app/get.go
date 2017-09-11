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

package app

import (
	"fmt"
	a "github.com/lastbackend/lastbackend/pkg/api/app/views/v1"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
)

func GetCmd(name string) {

	ns, err := Get(name)
	if err != nil {
		fmt.Print(err)
		return
	}

	ns.DrawTable()
}

func Get(name string) (*a.App, error) {

	var (
		err  error
		http = c.Get().GetHttpClient()
		er   = new(errors.Http)
		app  = new(a.App)
	)

	if len(name) == 0 {
		return nil, errors.BadParameter("name").Err()
	}

	_, _, err = http.
		GET(fmt.Sprintf("/app/%s", name)).
		AddHeader("Content-Type", "application/json").
		Request(&app, er)
	if err != nil {
		return nil, errors.New(er.Message)
	}

	if er.Code == 401 {
		return nil, errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, errors.New(er.Message)
	}

	return app, nil
}
