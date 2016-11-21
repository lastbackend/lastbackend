package project

import (
	"github.com/boltdb/bolt"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/libs/table"
)

func List() {

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
	res := []model.Project{}

	ctx.HTTP.
		GET("/project").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+token).
		Request(&res, &er)

	var header []string = []string{"ID", "Name", "Created", "Updated"}
	var data [][]string

	for i := 0; i < len(res); i++ {
		d := []string{
			res[i].ID,
			res[i].Name,
			res[i].Created.String()[:10],
			res[i].Updated.String()[:10],
		}

		data = append(data, d)
	}

	table.PrintTable(header, data, []string{})
}
