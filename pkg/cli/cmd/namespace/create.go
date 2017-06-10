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

package namespace

import (
	"fmt"
	n "github.com/lastbackend/lastbackend/pkg/api/namespace/views/v1"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
)

type createS struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

func CreateCmd(name, description string) {

	err := Create(name, description)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fmt.Sprintf("Namespace `%s` is created", name))
}

func Create(name, description string) error {

	var (
		err  error
		http = c.Get().GetHttpClient()
		er   = new(errors.Http)
		ns   = new(n.Namespace)
	)

	if len(name) == 0 {
		return errors.BadParameter("name").Err()
	}

	_, _, err = http.
		POST("/namespace").
		AddHeader("Content-Type", "application/json").
		BodyJSON(createS{name, description}).
		Request(&ns, er)
	if err != nil {
		return errors.New(er.Message)
	}

	if er.Code == 401 {
		return errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	ns, err = Switch(name)
	if err != nil {
		return err
	}

	return nil
}
