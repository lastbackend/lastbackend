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
	c "github.com/lastbackend/lastbackend/pkg/client/context"
	n "github.com/lastbackend/lastbackend/pkg/daemon/namespace/views/v1"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"log"
)

type createS struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

func CreateCmd(name, description string) {

	err := Create(name, description)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print("Namespace `" + name + "` is created")
}

func Create(name, description string) error {

	var (
		err       error
		http      = c.Get().GetHttpClient()
		er        = new(errors.Http)
		namespace = new(n.Namespace)
	)

	if len(name) == 0 {
		return errors.BadParameter("name").Err()
	}

	_, _, err = http.
		POST("/namespace").
		AddHeader("Content-Type", "application/json").
		BodyJSON(createS{name, description}).
		Request(&namespace, er)
	if err != nil {
		log.Println(err)
		return errors.New(er.Message)
	}

	if er.Code == 401 {
		return errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	namespace, err = Switch(name)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
