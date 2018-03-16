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
	"github.com/lastbackend/lastbackend/pkg/util/url"
)

type RequestServiceSpecS struct {
	Memory int64 `json:"memory,omitempty"`
}

type createS struct {
	Name    string              `json:"name"`
	Sources string              `json:"sources"`
	Spec    RequestServiceSpecS `json:"spec"`
}

func CreateCmd(name, sources string, memory int64) {

	err := Create(name, sources, memory)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Waiting for start service
	// TODO: Show spinner

	fmt.Println("Service `" + name + "` is succesfully created")
}

func Create(name, sources string, memory int64) error {

	var (
		err  error
		http = c.Get().GetHttpClient()
		er   = new(e.Http)
		res  = new(struct{})
	)

	s := url.Decode(sources)
	if (s.Hub == "index.docker.io" || s.Hub == "hub.lstbknd.net") && (s.Owner == "" || s.Name == "" || s.Branch == "") {
		return e.New("Incorrect sources")
	} else if s.Hub == "" || s.Owner == "" || s.Name == "" {
		return e.New("Incorrect sources")
	}

	ns, err := n.Current()
	if err != nil {
		return err
	}
	if ns.Meta == nil {
		return e.New("Workspace didn't select")
	}

	_, _, err = http.
		POST(fmt.Sprintf("/namespace/%s/service", ns.Meta.Name)).
		AddHeader("Content-Type", "application/json").
		BodyJSON(createS{Name: name, Sources: sources, Spec: RequestServiceSpecS{Memory: memory}}).
		Request(res, er)
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
