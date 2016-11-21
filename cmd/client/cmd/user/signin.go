package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/jarcoal/httpmock"
	mock "github.com/lastbackend/lastbackend/cmd/client/cmd/user/mocks"
	structs "github.com/lastbackend/lastbackend/cmd/client/cmd/user/structs"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	"github.com/lastbackend/lastbackend/libs/log/filesystem"
	"io/ioutil"
)

func SignIn(ctx *context.Context) {
	token, err, _ := Login(ctx)
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

func Login(ctx *context.Context) (string, error, string) {
	var password string
	var login string

	if ctx == context.Mock() {
		if ctx.Info.Version == "OK" {
			login, password = mock.MockSignInOk()
		} else if ctx.Info.Version == "BAD" {
			login, password = mock.MockSignInBad()
		}
		defer httpmock.Deactivate()
	} else {
		fmt.Print("Login: ")
		fmt.Scan(&login)

		fmt.Print("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			return "", err, ""
		}
		password = string(pass)
	}

	data := structs.LoginInfo{login, password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err, ""
	}

	resp, status := httpClient.Post(config.Get().AuthUserUrl, jsonData,
		"Content-Type", "application/json")
	if status == 200 {
		var token structs.TokenInfo
		err = json.Unmarshal(resp, &token)
		if err != nil {
			return "", err, ""
		}

		fmt.Println("Login successful")

		return token.Token, err, ""
	}

	var httpError structs.ErrorJson
	err = json.Unmarshal(resp, &httpError)
	if err != nil {
		return "", err, ""
	}
	fmt.Printf("Login failed: %s", httpError.Message)

	return "", nil, httpError.Message
}
