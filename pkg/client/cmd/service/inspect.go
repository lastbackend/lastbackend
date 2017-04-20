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
	nspace "github.com/lastbackend/lastbackend/pkg/client/cmd/namespace"
	c "github.com/lastbackend/lastbackend/pkg/client/context"
	s "github.com/lastbackend/lastbackend/pkg/daemon/service/views/v1"
	"github.com/lastbackend/lastbackend/pkg/errors"
)

func InspectCmd(name string) {

	service, namespace, err := Inspect(name)
	if err != nil {
		fmt.Print(err)
		return
	}

	service.DrawTable(namespace)
}

func Inspect(name string) (*s.Service, string, error) {

	var (
		err     error
		http    = c.Get().GetHttpClient()
		er      = new(errors.Http)
		service *s.Service
	)

	namespace, err := nspace.Current()
	if err != nil {
		return nil, "", errors.New(err.Error())
	}
	if namespace == nil {
		return nil, "", errors.New("Namespace didn't select")
	}

	_, _, err = http.
		GET("/namespace/"+namespace.Meta.Name+"/service/"+name).
		Request(&service, er)
	if err != nil {
		return nil, "", errors.New(er.Message)
	}

	if er.Code == 401 {
		return nil, "", errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, "", errors.New(er.Message)
	}

	return service, namespace.Meta.Name, nil
}
