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

package service

import (
	"fmt"
	s "github.com/lastbackend/lastbackend/pkg/api/service/views/v1"
	a "github.com/lastbackend/lastbackend/pkg/cli/cmd/app"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
)

func InspectCmd(name string) {

	srv, ns, err := Inspect(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	srv.DrawTable(ns)
}

func Inspect(name string) (*s.Service, string, error) {

	var (
		err  error
		http = c.Get().GetHttpClient()
		er   = new(errors.Http)
		srv  *s.Service
	)

	a, err := a.Current()
	if err != nil {
		return nil, "", err
	}
	if a == nil {
		return nil, "", errors.New("App didn't select")
	}

	_, _, err = http.
		GET(fmt.Sprintf("/app/%s/service/%s", a.Meta.Name, name)).
		Request(&srv, er)
	if err != nil {
		return nil, "", errors.New(er.Message)
	}

	if er.Code == 401 {
		return nil, "", errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, "", errors.New(er.Message)
	}

	return srv, a.Meta.Name, nil
}
