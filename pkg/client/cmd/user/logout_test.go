package user_test

import (
	f "github.com/lastbackend/lastbackend/utils"
	"io/ioutil"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/user"
)

func TestLogout(t *testing.T) {
	err := f.MkDir(f.GetHomeDir() + "/.lb", 0777)
	if err != nil {
		t.Error(err)
	}

	files, err := ioutil.ReadDir(f.GetHomeDir() + "/.lb")
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, files)

	err = user.Logout()
	if err != nil {
		t.Error(err)
	}

	files, err = ioutil.ReadDir(f.GetHomeDir() + "/.lb")
	if err != nil {
		assert.Nil(t, files)
	}
}
