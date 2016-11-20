package cmd

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/jarcoal/httpmock"
	mock "github.com/lastbackend/lastbackend/cmd/client/cmd/user/mocks"
	structs "github.com/lastbackend/lastbackend/cmd/client/cmd/user/structs"
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
		username, email, password = mock.MockSignUp()
		defer httpmock.Deactivate()
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

	data := structs.NewUserInfo{username, email, password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp, status := httpClient.Post(config.Get().UserUrl, jsonData, "Content-Type", "application/json")
	if status == 200 {
		fmt.Println("Account create")
	}

	var token structs.TokenInfo
	err = json.Unmarshal(resp, &token)
	if err != nil {
		return "", err
	}

	return token.Token, err
}
