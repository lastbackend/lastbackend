package project_test

import (
	"github.com/lastbackend/lastbackend/libs/db"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"testing"
	"time"
)

func TestCurrent(t *testing.T) {

	const token = "mocktoken"

	var (
		err  error
		ctx  = context.Mock()
		data = model.Project{
			Name:        "mock_name",
			ID:          "mock_id",
			Created:     time.Now(),
			Updated:     time.Now(),
			User:        "mock_demo",
			Description: "sample description",
		}
	)

	ctx.Storage, err = db.Init()
	if err != nil {
		t.Error(err)
		return
	}
	defer (func() {
		err = ctx.Storage.Clear()
		if err != nil {
			t.Error(err)
			return
		}
		err = ctx.Storage.Close()
		if err != nil {
			t.Error(err)
			return
		}
	})()

	ctx.Token = token

	if err != nil {
		t.Error(err)
		return
	}
	defer ctx.Storage.Close()

	err = ctx.Storage.Set("project", data)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = project.Current()
	if err != nil {
		t.Error(err)
		return
	}
}
