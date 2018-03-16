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

	cl "github.com/lastbackend/lastbackend/pkg/cli/cmd/cluster"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	v "github.com/lastbackend/lastbackend/pkg/cli/view"
	e "github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

type createS struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Cluster     string `json:"cluster"`
}

func CreateCmd(name, desc string) {

	err := Create(name, desc)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fmt.Sprintf("Workspace `%s` is created", name))
}

func Create(name, desc string) error {

	var (
		err      error
		http     = c.Get().GetHttpClient()
		er       = new(e.Http)
		clusters = new(v.ClusterList)
		response *v.Namespace
		storage  = c.Get().GetStorage()
	)

	clusters, err = cl.List()
	if err != nil {
		return err
	}

	_, _, err = http.
		POST("/namespace").
		AddHeader("Content-Type", "application/json").
		BodyJSON(createS{name, desc, (*clusters)[0].ID}).
		Request(&response, er)
	if err != nil {
		return e.New(er.Message)
	}

	if er.Code == 401 {
		return e.NotLoggedMessage
	}

	if er.Code == 404 {
		return e.New("Not found")
	}

	if er.Code == 500 {
		return e.UnknownMessage
	}

	if er.Code != 0 {
		return e.New(er.Message)
	}

	if err := storage.Namespace().Save(response); err != nil {
		return e.UnknownMessage
	}

	return nil
}
