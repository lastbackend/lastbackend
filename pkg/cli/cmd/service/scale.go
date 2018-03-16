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

	ns "github.com/lastbackend/lastbackend/pkg/cli/cmd/namespace"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	e "github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

type ScaleS struct {
	Replicas int64 `json:"replicas"`
}

func ScaleCmd(name string, replicas int64) {

	err := Scale(name, replicas)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Service `" + name + "` is succesfully scaled")
}

func Scale(name string, replicas int64) error {

	var (
		err      error
		http     = c.Get().GetHttpClient()
		er       = new(e.Http)
		response = new(struct{})
	)

	n, err := ns.Current()
	if err != nil {
		return err
	}
	if n.Meta == nil {
		return e.New("Workspace didn't select")
	}

	s, err := Inspect(name)
	if err != nil {
		return err
	}
	if len(s.Deployments) == 0 {
		return e.New("Service now is not deployed")
	}

	var deployment = ""
	for _, d := range s.Deployments {
		if d.State.Active {
			deployment = d.ID
		}
	}

	_, _, err = http.
		PUT(fmt.Sprintf("/namespace/%s/service/%s/deployment/%s", n.Meta.Name, s.Meta.Name, deployment)).
		AddHeader("Content-Type", "application/json").
		BodyJSON(ScaleS{Replicas: replicas}).
		Request(response, er)
	if err != nil {
		return e.New(er.Message)
	}

	if er.Code == 401 {
		return e.NotLoggedMessage
	}

	if er.Code != 0 {
		return e.New(er.Message)
	}

	return nil
}
