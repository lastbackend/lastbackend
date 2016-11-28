package service

import (
tab "github.com/crackcomm/go-clitable"
e "github.com/lastbackend/lastbackend/libs/errors"
"github.com/lastbackend/lastbackend/libs/model"
"github.com/lastbackend/lastbackend/pkg/client/context"
)

type serviceCreate struct {
	Name string `json:"name"`
}

func getToken(ctx *context.Context) (string, error) {
	var err error
	token := struct {
		Token string `json:"token"`
	}{}
	err = ctx.Storage.Get("session", &token)
	return token.Token, err
}

func printData(data model.Service) {
	table := tab.New([]string{"ID", "Project", "Name", "Created", "Updated"})
	table.AddRow(map[string]interface{}{
		"ID":      data.ID,
		"Project": data.Image,
		"Name":    data.Name,
		"Created": data.Created.String()[:10],
		"Updated": data.Updated.String()[:10],
	})
	table.Markdown = true
	table.Print()
}

func CreateCmd(name string) {

	ctx := context.Get()

	err := Create(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}
}

func Create(name string) error {
	var (
		err   error
		ctx   = context.Get()
		token string
		res   model.Project
	)
	token, err = getToken(ctx)
	req_err := new(e.Http)
	_, _, err = ctx.HTTP.
		POST("/service").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+token).
		BodyJSON(serviceCreate{name}).
		Request(&res, req_err)
	if err != nil {
	}
	return err

}



