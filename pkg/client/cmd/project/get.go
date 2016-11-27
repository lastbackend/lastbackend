package project

import (
	"errors"
	tab "github.com/crackcomm/go-clitable"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func GetCmd(name string) {

	var ctx = context.Get()

	err := Get(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}
}

func Get(name string) error {

	var (
		err error
		ctx = context.Get()
	)

	token := struct {
		Token string `json:"token"`
	}{}

	err = ctx.Storage.Get("session", &token)
	if err != nil {
		return errors.New(err.Error())
	}
	if token.Token == "" {
		return errors.New(e.StatusAccessDenied)
	}

	er := new(e.Http)
	res := new(model.Project)

	_, _, err = ctx.HTTP.
		GET("/project/"+name).
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+token.Token).
		Request(&res, er)
	if err != nil {
		return errors.New(err.Error())
	}

	if er.Code != 0 {
		return errors.New(e.Message(er.Status))
	}

	table := tab.New([]string{"ID", "Name", "Created", "Updated"})
	table.AddRow(map[string]interface{}{
		"ID":      res.ID,
		"Name":    res.Name,
		"Created": res.Created.String()[:10],
		"Updated": res.Updated.String()[:10],
	})
	table.Markdown = true
	table.Print()

	return nil
}
