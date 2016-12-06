package user_test

import (
	"github.com/lastbackend/lastbackend/libs/db"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/user"
	"github.com/lastbackend/lastbackend/pkg/client/context"
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

	ctx.Storage, err = db.Init()
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
