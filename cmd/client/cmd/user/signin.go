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

func SignIn() {
	token, err, _ := Login()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if token == "" {
		return
	}

	byteToken := []byte(token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = filesystem.MkDir(cfg.StoragePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = ioutil.WriteFile(cfg.StoragePath+"token", byteToken, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func Login() (string, error, string) {
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
		fmt.Println("LOGININGO")
		return "", err, ""
	}

	ctx.HTTP.POST("/user", struct {}{})

	resp, status := httpClient.Post(cfg.AuthUserUrl, jsonData,
		"Content-Type", "application/json")
	if status == 200 {
		var token structs.TokenInfo
		err = json.Unmarshal(resp, &token)
		if err != nil {
			fmt.Println("RESP")
			return "", err, ""
		}

		fmt.Println("Login successful")

		return token.Token, err, ""
	}

	var httpError structs.ErrorJson
	err = json.Unmarshal(resp, &httpError)
	if err != nil {
		fmt.Println("ERRORJSON")
		return "", err, ""
	}
	fmt.Printf("Login failed: %s\n", httpError.Message)

	return "", nil, httpError.Message
}
