package cmd

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	"k8s.io/client-go/1.5/pkg/util/json"
)

func SignUp(ctx *context.Context) {
	_, err := CreateNewUser()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func CreateNewUser() (string, error) {
	fmt.Print("Username: ")
	var username string
	fmt.Scan(&username)

	fmt.Print("Email: ")
	var email string
	fmt.Scan(&email)

	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		return "", err
	}
	password := string(pass)

	data := newUserInfo{username, email, password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp := httpClient.Post(config.Get().CreateUserUrl, jsonData, "Content-Type", "application/json")

	var token tokenInfo
	err = json.Unmarshal(resp, &token)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return token.Token, err
}
