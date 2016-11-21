package user

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/howeyc/gopass"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/lastbackend/lastbackend/libs/errors"
)

type UserCreateS struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUp() {

	var (
		err      error
		ctx      = context.Get()
		username string
		email    string
		password string
	)

	username, email, password, err = inputSignupUserData()
	if err != nil {
		fmt.Println(err.Error())
	}

	er := errors.Http{}
	res := struct {
		Token string `json:"token"`
	}{}

	ctx.HTTP.
		POST("/user").
		AddHeader("Content-Type", "application/json").
		BodyJSON(UserCreateS{username, email, password}).
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

func instruction() {
	fmt.Println("User field must be at least 4 letters")
	fmt.Println("Password filed must be at least 6 letters")
	fmt.Println("---Example---")
	fmt.Println("User: user")
	fmt.Println("Email: email@email.email")
	fmt.Println("Password: password")
	fmt.Println("-------------")
}

func inputSignupUserData() (string, string, string, error) {
	instruction()

	var username string
	var email string
	var password string

	fmt.Print("Username: ")
	fmt.Scan(&username)

	fmt.Print("Email: ")
	fmt.Scan(&email)

	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		return "", "", "", err
	}
	password = string(pass)

	return username, email, password, err
}
