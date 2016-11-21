package project

import (
	"github.com/boltdb/bolt"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/libs/table"
)

func Get(name string) {

	var (
		ctx   = context.Get()
		token string
	)

	err := ctx.Storage.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("session"))
		buf := bucket.Get([]byte("token"))
		token = string(buf)

		return nil
	})
	if err != nil {
		ctx.Log.Error(err)
	}

	er := errors.Http{}
	res := model.Project{}

	ctx.HTTP.
		GET("/project/"+name).
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+token).
		Request(&res, &er)

	var header []string = []string{"ID", "Name", "Created", "Updated"}
	var data [][]string

	d := []string{
		res.ID,
		res.Name,
		res.Created.String()[:10],
		res.Updated.String()[:10],
	}

	data = append(data, d)
	d = d[:0]

	table.PrintTable(header, data, []string{})
}
