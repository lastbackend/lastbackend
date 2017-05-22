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
	"fmt"
	s "github.com/lastbackend/lastbackend/pkg/api/service/views/v1"
	n "github.com/lastbackend/lastbackend/pkg/cli/cmd/namespace"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
)

func ListServiceCmd() {

	srvList, ns, err := List()
	if err != nil {
		fmt.Print(err)
		return
	}

	if srvList != nil {
		srvList.DrawTable(ns)
	}
}

func List() (*s.ServiceList, string, error) {

	var (
		err     error
		http    = c.Get().GetHttpClient()
		er      = new(errors.Http)
		srvList *s.ServiceList
	)

	ns, err := n.Current()
	if err != nil {
		return nil, "", err
	}

	if ns == nil {
		return nil, "", errors.New("Namespace didn't select")
	}

	_, _, err = http.
		GET(fmt.Sprintf("/namespace/%s/service", ns.Meta.Name)).
		Request(&srvList, er)
	if err != nil {
		return nil, "", err
	}

	if er.Code == 401 {
		return nil, "", errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, "", errors.New(er.Message)
	}

	if len(*srvList) == 0 {
		return nil, "", errors.New("You don't have any services")
	}

	return srvList, ns.Meta.Name, nil
}
