package project_test

import (
	//"errors"
	//"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	//"github.com/kubernetes/client-go/pkg/apis/storage"
	"github.com/lastbackend/lastbackend/libs/db"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	//"time"
	"testing"
)
/*
var mock_proj = model.Project{Name:					"mock_name",
	ID:						"mock_id",
	Created:			time.Now(),
	Updated:			time.Now(),
	User: 				"you",
	Description:	"sample description"}
*/

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
	/*
	if ctx.Token == "" {
		return errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	}

	err = ctx.Storage.Get("project", project)
	if err != nil {
		return errors.New(err.Error())
	}

	if project.ID == "" {
		ctx.Log.Info("Project didn't select")
		return nil
	}

	project.DrawTable()

	return
	*/
}
