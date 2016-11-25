package user_test

import (
	"github.com/lastbackend/lastbackend/pkg/client/cmd/user"
	f "github.com/lastbackend/lastbackend/utils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestLogout(t *testing.T) {
	err := f.MkDir(f.GetHomeDir()+"/.lb", 0777)
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
