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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	c "github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"strings"
	"time"
)

type updateS struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

func UpdateCmd(name, newName, description string) {

	var (
		choice string
	)

	if description == "" {
		fmt.Println("Description is empty, field will be cleared\n" +
			"Want to continue? [Y\\n]")

		for {
			fmt.Scan(&choice)

			switch strings.ToLower(choice) {
			case "y":
				break
			case "n":
				return
			default:
				fmt.Print("Incorrect input. [Y\n]")
				continue
			}

			break
		}
	}

	err := Update(name, newName, description)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print("Namespace `" + name + "` is succesfully updated")
}

func Update(name, newName, description string) error {

	var (
		err     error
		http    = c.Get().GetHttpClient()
		storage = c.Get().GetStorage()
		er      = new(errors.Http)
		res     = new(types.Namespace)
	)

	_, _, err = http.
		PUT("/namespace/"+name).
		AddHeader("Content-Type", "application/json").
		BodyJSON(updateS{newName, description}).
		Request(&res, er)
	if err != nil {
		return errors.New(er.Message)
	}

	if er.Code == 401 {
		return errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	namespace, err := Current()
	if err != nil {
		return errors.New(err.Error())
	}

	if namespace != nil {
		if name == namespace.Meta.Name {
			namespace.Meta.Name = newName
			namespace.Meta.Description = description
			namespace.Meta.Updated = time.Now()

			if c.Get().IsMock() {
				err = storage.Set("test", namespace)
				if err != nil {
					return errors.UnknownMessage
				}
			} else {
				err = storage.Set("namespace", namespace)
				if err != nil {
					return errors.UnknownMessage
				}
			}
		}
	}

	return nil
}
