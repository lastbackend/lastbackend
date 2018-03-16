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

package namespace

import (
	"fmt"

	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	v "github.com/lastbackend/lastbackend/pkg/cli/view"
	e "github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

func CurrentCmd(name string) {

	if name == "" {
		ns, err := Current()
		if err != nil {
			fmt.Println(err)
			return
		}

		if ns.Meta.Name == "" {
			fmt.Print("Workspace didn't select")
			return
		}

		ns.Print()
	} else {
		ns, err := Get(name)
		if err != nil {
			fmt.Println(err)
			return
		}

		ns.Print()
	}
}

func Get(name string) (*v.Namespace, error) {

	var (
		err      error
		http     = c.Get().GetHttpClient()
		er       = new(e.Http)
		response = new(v.Namespace)
	)

	_, _, err = http.
		GET(fmt.Sprintf("/namespace/%s", name)).
		AddHeader("Content-Type", "application/json").
		Request(&response, er)
	if err != nil {
		return nil, e.UnknownMessage
	}

	if er.Code == 401 {
		return nil, e.NotLoggedMessage
	}

	if er.Code == 404 {
		return nil, e.New(er.Message)
	}

	if er.Code == 500 {
		return nil, e.UnknownMessage
	}

	if er.Code != 0 {
		return nil, e.New(er.Message)
	}

	return response, nil
}

func Current() (*v.Namespace, error) {

	var (
		err     error
		storage = c.Get().GetStorage()
	)

	ns, err := storage.Namespace().Load()
	if err != nil {
		return nil, e.UnknownMessage
	}

	return ns, nil
}
