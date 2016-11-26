package project

import (
	"errors"
	tab "github.com/crackcomm/go-clitable"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func ListCmd() {

	var ctx = context.Get()

	err := List()
	if err != nil {
		ctx.Log.Error(err)
		return
	}
}

func List() error {

	var (
		err   error
		ctx   = context.Get()
		token *string
	)

	err = ctx.Storage.Get("session", nil)
	if token == nil {
		return errors.New(e.StatusAccessDenied)
	}

	er := new(e.Http)
	res := []model.Project{}

	_, _, err = ctx.HTTP.
		GET("/project").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+*token).
		Request(&res, er)
	if err != nil {
		return err
	}

	if er.Code != 0 {
		return errors.New(e.Message(er.Status))
	}

	for i := 0; i < len(res); i++ {
		table := tab.New([]string{"ID", "NAME", "Created", "Updated"})
		table.AddRow(map[string]interface{}{
			"ID":      res[i].ID,
			"Name":    res[i].Name,
			"Created": res[i].Created.String()[:10],
			"Updated": res[i].Updated.String()[:10],
		})
		table.Markdown = true
		table.Print()
	}

	return nil
}
