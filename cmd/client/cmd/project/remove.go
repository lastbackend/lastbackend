package project

import (
	"github.com/boltdb/bolt"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/lastbackend/lastbackend/libs/errors"
)

func Remove(name string) {

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
	res := struct{}{}

	ctx.HTTP.
		DELETE("/project/"+name).
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+token).
		Request(&res, &er)
}
