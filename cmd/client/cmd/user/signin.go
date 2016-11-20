package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/jarcoal/httpmock"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	structs "github.com/lastbackend/lastbackend/cmd/client/cmd/user/structs"
	mock "github.com/lastbackend/lastbackend/cmd/client/cmd/user/mocks"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	"github.com/lastbackend/lastbackend/libs/log/filesystem"
	"io/ioutil"
)

func SignIn(ctx *context.Context) {
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

	err = filesystem.MkDir(config.Get().StoragePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = ioutil.WriteFile(config.Get().StoragePath+"token", byteToken, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func Login(ctx *context.Context) (string, error) {
	var password string
	var login string

	if ctx == context.Mock() {
		login, password = mock.MockAuth()
		defer httpmock.Deactivate()
	} else {
		fmt.Print("Login: ")
		fmt.Scan(&login)

		fmt.Print("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			return "", err
		}
		password = string(pass)
	}

	data := structs.LoginInfo{login, password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp := httpClient.Post(config.Get().AuthUserUrl, jsonData, "Content-Type", "application/json")

	var token structs.TokenInfo
	err = json.Unmarshal(resp, &token)
	if err != nil {
		return "", err
	}

	return token.Token, err
}
