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

package project

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

func UpdateCmd(name, newProjectName, description string) {

	var (
		log    = c.Get().GetLogger()
		choice string
	)

	if description == "" {
		log.Info("Description is empty, field will be cleared\n" +
			"Want to continue? [Y\\n]")

		for {
			fmt.Scan(&choice)

			switch strings.ToLower(choice) {
			case "y":
				break
			case "n":
				return
			default:
				log.Error("Incorrect input. [Y\n]")
				continue
			}

			break
		}
	}

	err := Update(name, newProjectName, description)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Successful")
}

func Update(name, newProjectName, description string) error {

	var (
		err     error
		http    = c.Get().GetHttpClient()
		storage = c.Get().GetStorage()
		er      = new(errors.Http)
		res     = new(types.Namespace)
	)

	_, _, err = http.
		PUT("/project/"+name).
		AddHeader("Content-Type", "application/json").
		BodyJSON(updateS{newProjectName, description}).
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

	project, err := Current()
	if err != nil {
		return errors.New(err.Error())
	}

	if project != nil {
		if name == project.Name {
			project.Name = newProjectName
			project.Description = description
			project.Updated = time.Now()

			err = storage.Set("project", project)
			if err != nil {
				return errors.New(err.Error())
			}
		}
	}

	return nil
}
