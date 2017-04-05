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

package user_test

import (
	"github.com/lastbackend/lastbackend/pkg/client/cmd/user"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/lastbackend/lastbackend/pkg/client/storage"
	"github.com/lastbackend/lastbackend/pkg/util/homedir"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestLogout(t *testing.T) {

	var (
		err error
		ctx = context.Mock()
	)

	ctx.Storage, err = storage.Init()
	if err != nil {
		panic(err)
	}
	defer (func() {
		err = ctx.Storage.Close()
		if err != nil {
			t.Error(err)
			return
		}
	})()

	err = user.Logout()
	if err != nil {
		t.Error(err)
	}

	files, err := ioutil.ReadDir(homedir.HomeDir() + "/.lb")
	if err != nil {
		assert.Nil(t, files)
	}
}
