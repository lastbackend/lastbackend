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

package namespace

import (
	"fmt"
	c "github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
)

func RemoveCmd(name string) {

	if err := Remove(name); err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print("namespace `" + name + "` is successfully removed")
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
		DELETE("/namespace/"+name).
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

	namespace, err := Current()
	if err != nil {
		return errors.New(err.Error())
	}

	if namespace != nil {
		if name == namespace.Meta.Name {
			var sName string
			if c.Get().IsMock() {
				sName = "test"
			} else {
				sName = "namespace"
			}
			err = storage.Set(sName, nil)
			if err != nil {
				return errors.UnknownMessage
			}
		}
	}

	return nil
}
