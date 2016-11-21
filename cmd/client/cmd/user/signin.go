package user

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/howeyc/gopass"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/lastbackend/lastbackend/libs/errors"
)

type UserLoginS struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func SignIn() {

	var (
		err      error
		login    string
		password string
		ctx      = context.Get()
	)

	login, password, err = inputLoginUserData()
	if err != nil {
		fmt.Println(err.Error())
	}

	er := errors.Http{}
	res := struct {
		Token string `json:"token"`
	}{}

	ctx.HTTP.
		POST("/session").
		AddHeader("Content-Type", "application/json").
		BodyJSON(UserLoginS{login, password}).
		Request(&res, &er)

	if res.Token == "" {
		return
	}

	err = ctx.Storage.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("session"))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte("token"), []byte(res.Token))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		ctx.Log.Error(err)
	}
}

func inputLoginUserData() (string, string, error) {

	var (
		login    string
		password string
	)

	fmt.Print("Login: ")
	fmt.Scan(&login)

	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		return "", "", err
	}
	password = string(pass)

	return login, password, err
}
