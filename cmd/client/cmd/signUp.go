package cmd

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/jarcoal/httpmock"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	"k8s.io/client-go/1.5/pkg/util/json"
)

func SignUp(ctx *context.Context) {
	_, err := CreateNewUser(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func CreateNewUser(ctx *context.Context) (string, error) {
	var username string
	var email string
	var password string

	if ctx == context.Mock() {
		username, email, password = MockSignUp()
	} else {
		fmt.Print("Username: ")
		fmt.Scan(&username)

		fmt.Print("Email: ")
		fmt.Scan(&email)

		fmt.Print("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			return "", err
		}
		password = string(pass)
	}

	data := newUserInfo{username, email, password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp := httpClient.Post(config.Get().UserUrl, jsonData, "Content-Type", "application/json")

	var token tokenInfo
	err = json.Unmarshal(resp, &token)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return token.Token, err
}

func MockSignUp() (string, string, string) {
	username := "testname"
	email := "test@lb.com"
	password := "12345678"

	httpmock.Activate()

	httpmock.RegisterResponder("POST", config.Get().UserUrl,
		httpmock.NewStringResponder(200, `{"token": "token"}`))

	return username, email, password
}
