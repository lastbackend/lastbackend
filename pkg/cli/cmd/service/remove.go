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
	nspace "github.com/lastbackend/lastbackend/pkg/cli/cmd/app"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"github.com/lastbackend/lastbackend/pkg/common/types"
)

func RemoveCmd(name string) {

	err := Remove(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Service `" + name + "` succesfully removed")
}

func Remove(name string) error {

	var (
		err  error
		http = c.Get().GetHttpClient()
		srv  = new(types.App)
		er   = new(errors.Http)
	)

	ns, err := nspace.Current()
	if err != nil {
		return err
	}

	_, _, err = http.
		DELETE(fmt.Sprintf("/app/%s/service/%s", ns.Meta.Name, name)).
		Request(srv, er)
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
