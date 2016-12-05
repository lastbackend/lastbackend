package project_test

import (
	"github.com/lastbackend/lastbackend/libs/db"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"testing"
)

func TestCurrent(t *testing.T) {
	var err error
	ctx := context.Mock()

	ctx.Storage, err = db.Init()

	if err != nil {
		t.Error(err)
		return
	}
	defer ctx.Storage.Close()
	ctx.Token = token

	err = ctx.Storage.Set("project", mock_proj)

	if err != nil {
		t.Error(err)
		return
	}

	err = project.Current()

	if err != nil {
		t.Error(err)
		return
	}

}
