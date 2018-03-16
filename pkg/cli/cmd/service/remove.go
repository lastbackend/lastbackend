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

	n "github.com/lastbackend/lastbackend/pkg/cli/cmd/namespace"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	e "github.com/lastbackend/lastbackend/pkg/distribution/errors"
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
		res  = new(struct{})
		er   = new(e.Http)
	)

	ns, err := n.Current()
	if err != nil {
		return err
	}
	if ns.Meta == nil {
		return e.New("Workspace didn't select")
	}

	_, _, err = http.
		DELETE(fmt.Sprintf("/namespace/%s/service/%s", ns.Meta.Name, name)).
		AddHeader("Content-Type", "application/json").
		Request(&res, er)
	if err != nil {
		return e.UnknownMessage
	}

	if er.Code == 401 {
		return e.NotLoggedMessage
	}

	if er.Code != 0 {
		return e.New(er.Message)
	}

	return nil
}
