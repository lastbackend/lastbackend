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
	n "github.com/lastbackend/lastbackend/pkg/api/app/views/v1"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
)

func SwitchCmd(name string) {

	ns, err := Switch(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The app `%s` was selected as the current\n", ns.Meta.Name)
}

func Switch(name string) (*n.App, error) {

	var (
		er      = new(errors.Http)
		http    = c.Get().GetHttpClient()
		storage = c.Get().GetStorage()
		a       = new(n.App)
	)

	_, _, err := http.
		GET(fmt.Sprintf("/app/%s", name)).
		AddHeader("Content-Type", "application/json").
		Request(&a, er)
	if err != nil {
		return nil, errors.New(er.Message)
	}

	if er.Code == 401 {
		return nil, errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, errors.New(er.Message)
	}

	if err := storage.App().Save(a); err != nil {
		return nil, errors.UnknownMessage
	}

	return a, nil
}
