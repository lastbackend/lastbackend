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

	var mock_proj = model.Project{Name: "mock_name",
		ID:          "mock_id",
		Created:     time.Now(),
		Updated:     time.Now(),
		User:        "you",
		Description: "sample description"}

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
