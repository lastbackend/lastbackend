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

package user

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	e "github.com/lastbackend/lastbackend/pkg/errors"
)

func LogoutCmd() {

	ctx := context.Get()

	err := Logout()
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	ctx.Log.Info("Logout successfully")
	fmt.Println("Logout successfully")
}

func Logout() error {

	var (
		err error
		ctx = context.Get()
	)

	err = ctx.Storage.Clear()
	if err != nil {
		return e.LogoutErrorMessage
	}

	return nil
}
