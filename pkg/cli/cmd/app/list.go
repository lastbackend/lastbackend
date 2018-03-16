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

package app

import (
	"fmt"
	a "github.com/lastbackend/lastbackend/pkg/api/app/views/v1"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
)

func ListAppCmd() {

	nsList, err := List()
	if err != nil {
		fmt.Println(err)
		return
	}

	if nsList != nil {
		nsList.DrawTable()
	}
}

func List() (*a.AppList, error) {

	var (
		err     error
		http    = c.Get().GetHttpClient()
		er      = new(errors.Http)
		appList = new(a.AppList)
	)

	_, _, err = http.
		GET("/app").
		Request(appList, er)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if er.Code == 401 {
		return nil, errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, errors.New(er.Message)
	}

	if len(*appList) == 0 {
		return nil, errors.New("You don't have any app")
	}

	return appList, nil
}
