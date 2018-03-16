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

package cluster

import (
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	v "github.com/lastbackend/lastbackend/pkg/cli/view"
	e "github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

func List() (*v.ClusterList, error) {

	var (
		err      error
		http     = c.Get().GetHttpClient()
		er       = new(e.Http)
		response = new(v.ClusterList)
	)

	_, _, err = http.
		GET("/cluster").
		AddHeader("Content-Type", "application/json").
		Request(response, er)
	if err != nil {
		return nil, e.UnknownMessage
	}

	if er.Code == 401 {
		return nil, e.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, e.New(er.Message)
	}

	return response, nil
}
