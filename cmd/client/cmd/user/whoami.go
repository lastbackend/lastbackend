package user

import (
	"github.com/boltdb/bolt"
	"github.com/lastbackend/api/libs/model"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/table"
	"strconv"
)

func Whoami() {

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
	res := model.User{}

	ctx.HTTP.
		GET("/user").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+token).
		Request(&res, &er)

	var header []string = []string{"Username", "Email", "Balance", "Organization", "Created", "Updated"}
	var data [][]string

	organization := strconv.FormatBool(res.Organization)
	balance := strconv.FormatFloat(float64(res.Balance), 'f', 2, 32)
	d := []string{
		res.Username,
		res.Email,
		balance,
		organization,
		res.Created.String()[:10],
		res.Updated.String()[:10],
	}

	data = append(data, d)
	d = d[:0]

	table.PrintTable(header, data, []string{})
}
