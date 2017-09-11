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

package app

import (
	"fmt"
	a "github.com/lastbackend/lastbackend/pkg/api/app/views/v1"
	c "github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
)

func CurrentCmd() {

	app, err := Current()
	if err != nil {
		fmt.Print(err)
		return
	}

	if app == nil {
		fmt.Print("App didn't select")
		return
	}

	app.DrawTable()
}

func Current() (*a.App, error) {

	var (
		err     error
		storage = c.Get().GetStorage()
	)

	app, err := storage.App().Load()
	if err != nil {
		return nil, errors.UnknownMessage
	}

	return app, nil
}
