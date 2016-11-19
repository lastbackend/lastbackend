package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/jarcoal/httpmock"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	"io/ioutil"
)

func Auth(ctx *context.Context) {
	token, err := Login(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	byteToken := []byte(token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = ioutil.WriteFile(config.Get().StoragePath, byteToken, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func Login(ctx *context.Context) (string, error) {
	var password string
	var login string

	if ctx == context.Mock() {
		login, password = MockUp()
		defer httpmock.Deactivate()
	} else {
		fmt.Print("Login: ")
		fmt.Scan(&login)

		fmt.Print("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			fmt.Println(err.Error())
			return "", err
		}
		password = string(pass)
	}

	data := loginInfo{login, password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	resp := httpClient.Post(config.Get().AuthUserUrl, jsonData, "Content-Type", "application/json")

	var token tokenInfo
	err = json.Unmarshal(resp, &token)
	if err != nil {
		fmt.Println(err.Error())
	}

	return token.Token, err
}

func MockUp() (string, string) {
	login := "lavr"
	password := "12345678"

	httpmock.Activate()

	httpmock.RegisterResponder("POST", config.Get().AuthUserUrl,
		httpmock.NewStringResponder(200, `{"token": "token"}`))

	return login, password
}
