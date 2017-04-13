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

func UpdateCmd(name, newNamespace, description string) {

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

	err := Update(name, newNamespace, description)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print("Successful")
}

func Update(name, newNamespace, description string) error {

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
		BodyJSON(updateS{newNamespace, description}).
		Request(&res, er)
	if err != nil {
		return err
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
		if name == namespace.Name {
			namespace.Name = newNamespace
			namespace.Description = description
			namespace.Updated = time.Now()

			err = storage.Set("namespace", namespace)
			if err != nil {
				return errors.New(err.Error())
			}
		}
	}

	return nil
}
