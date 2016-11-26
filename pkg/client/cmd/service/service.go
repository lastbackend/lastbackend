package service

import (
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	tab "github.com/crackcomm/go-clitable"
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

func Create(name string) error {
	var (
		err error
		ctx = context.Get()
		token string
		res model.Project
	)
	token, err = getToken(ctx)
	req_err := new(e.Http)
	_, _, err = ctx.HTTP.
		POST("/service").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer " + token).
		BodyJSON(serviceCreate{name}).
		Request(&res, req_err)
	if err != nil {}
	return err

}

func Inspect(name string) error {
	var (
		err error
		ctx = context.Get()
		token string
		res model.Service
	)
	token, err = getToken(ctx)
	req_err := new(e.Http)
	_, _, err = ctx.HTTP.
		GET("/service/" + name).
		AddHeader("Authorization", "Bearer " + token).
		Request(&res, req_err)
	if err != nil {
		return err
	}


	printData(res)
	return err
}

//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbSI6Im1vY2tlZEBtb2NrZWQuY29tIiwiZXhwIjoxNDg3NDE3ODk1LCJqdGkiOjE0ODc0MTc4OTUsIm9pZCI6IiIsInVpZCI6IjU2MmYwY2EwLTI2ZWEtNGFiNC1hZDBmLTU1N2NmYjJmYjgwNyIsInVzZXIiOiJtb2NrZWQifQ.VjHgKRqJCwf7TDphHPHhMl6njwL7agE1dzPVeGy5HFI

func Remove(name string) error {
	var (
		err error
		ctx = context.Get()
		token string
		res model.Project
	)
	token, err = getToken(ctx)
	req_err := new(e.Http)
	_, _, err = ctx.HTTP.
		DELETE("/service/" + name).
		AddHeader("Authorization", "Bearer " + token).
		Request(&res, req_err)
	if err != nil {}
	return err
}