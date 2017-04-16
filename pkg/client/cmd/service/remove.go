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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	nspace "github.com/lastbackend/lastbackend/pkg/client/cmd/namespace"
	c "github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
)

func RemoveCmd(name string) {

	err := Remove(name)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print("Successful")
}

func Remove(name string) error {

	var (
		err     error
		http    = c.Get().GetHttpClient()
		service = new(types.Namespace)
		er      = new(errors.Http)
	)

	namespace, err := nspace.Current()
	if err != nil {
		return errors.New(err.Error())
	}

	_, _, err = http.
		DELETE("/namespace/"+namespace.Name+"/service/"+name).
		Request(service, er)
	if err != nil {
		return errors.New(err.Error())
	}

	if er.Code == 401 {
		return errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	return nil
}
