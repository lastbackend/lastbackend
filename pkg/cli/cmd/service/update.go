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
	v "github.com/lastbackend/lastbackend/pkg/cli/view"
	e "github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

type updateS struct {
	Name string      `json:"name"`
	Spec OptionsSpec `json:"spec"`
}

type OptionsSpec struct {
	Replicas int64 `json:"replicas,omitempty"`
	Memory   int64 `json:"memory,omitempty"`
}

func UpdateCmd(name string, memory int64) {

	err := Update(name, memory)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Service `" + name + "` successfully updated")
}

func Update(name string, memory int64) error {

	var (
		err      error
		http     = c.Get().GetHttpClient()
		er       = new(e.Http)
		response = new(v.Service)
	)

	srv, err := Inspect(name)
	if err != nil {
		return err
	}

	var cfg = updateS{}

	if memory != 0 {
		cfg.Spec.Memory = memory
	} else {
		cfg.Spec.Memory = srv.Spec.Memory
	}

	cfg.Name = srv.Meta.Name

	ns, err := n.Current()
	if err != nil {
		return err
	}
	if ns.Meta == nil {
		return e.New("Workspace didn't select")
	}

	_, _, err = http.
		PUT(fmt.Sprintf("/namespace/%s/service/%s", ns.Meta.Name, name)).
		AddHeader("Content-Type", "application/json").
		BodyJSON(cfg).
		Request(&response, er)
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
