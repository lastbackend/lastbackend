package user_test

import (
	"github.com/lastbackend/lastbackend/pkg/client/cmd/user"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	f "github.com/lastbackend/lastbackend/utils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestLogout(t *testing.T) {
	_ = context.Get()

	var (
		err error
		ctx = context.Mock()
	)

	err = ctx.Storage.Init()
	if err != nil {
		panic(err)
	}
	defer ctx.Storage.Close()

	err = user.Logout()
	if err != nil {
		t.Error(err)
	}

	files, err := ioutil.ReadDir(f.GetHomeDir() + "/.lb")
	if err != nil {
		assert.Nil(t, files)
	}
}
