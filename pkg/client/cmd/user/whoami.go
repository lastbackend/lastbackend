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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
)

func WhoamiCmd() {

	var (
		err error
		ctx = context.Get()
	)

	user, err := Whoami()
	if err != nil {
		fmt.Println(err)
		ctx.Log.Error(err)
		return
	}

	ctx.Log.Info(fmt.Sprintf("User information:\r\n\r\n"+
		"Username: \t%s\n"+
		"E-mail: \t%s\n"+
		"Created: \t%s\n"+
		"Updated: \t%s", user.Username, user.Emails.GetDefault(),
		user.Created.String()[:10], user.Updated.String()[:10]))
}

func Whoami() (*types.User, error) {

	var (
		err  error
		ctx  = context.Get()
		er   = new(errors.Http)
		user = new(types.User)
	)

	if ctx.Token == "" {
		return nil, errors.NotLoggedMessage
	}

	_, _, err = ctx.HTTP.
		GET("/user").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(user, er)
	if err != nil {
		return nil, err
	}

	if er.Code == 401 {
		return nil, errors.NotLoggedMessage
	}

	if er.Code == 500 {
		return nil, errors.UnknownMessage
	}

	if er.Code != 0 {
		return nil, errors.New(er.Message)
	}

	return user, nil
}
