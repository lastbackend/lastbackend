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
	"github.com/lastbackend/lastbackend/libs/log/filesystem"
	"io/ioutil"
	"k8s.io/client-go/1.5/pkg/util/json"
)

func SignUp(ctx *context.Context, cfg *config.Config) {
	token, err, _ := CreateNewUser(ctx, cfg)
	if err != nil {
		fmt.Println(err.Error())
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

func instruction() {
	fmt.Println("User field must be at least 4 letters")
	fmt.Println("Email filed must be at least 6 letters")
	fmt.Println("Password filed must be at least 6 letters")
}

func inputUserData() (string, string, string) {
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
		return "", err, ""
	}
	password = string(pass)

	return username, email, password
}

func CreateNewUser(ctx *context.Context, cfg *config.Config) (string, error, string) {
	var username string
	var email string
	var password string

	if ctx == context.Mock() {
		if ctx.Info.Version == "OK" {
			username, email, password = mock.MockSignUpOk()
		} else if ctx.Info.Version == "BAD_USERNAME" {
			username, email, password = mock.MockSignUpBadUsername()
		} else if ctx.Info.Version == "BAD_EMAIL" {
			username, email, password = mock.MockSignUpBadEmail()
		} else if ctx.Info.Version == "BAD_PASSWORD" {
			username, email, password = mock.MockSignUpBadPassword()
		}
		defer httpmock.Deactivate()
	} else {
		username, email, password = inputUserData()
	}

	data := structs.NewUserInfo{username, email, password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err, ""
	}

	resp, status := httpClient.Post(cfg.UserUrl, jsonData, "Content-Type", "application/json")
	if status == 200 {
		var token structs.TokenInfo
		err = json.Unmarshal(resp, &token)
		if err != nil {
			return "", err, ""
		}
		fmt.Println("Account created successful")

		return token.Token, err, ""
	}

	var httpError structs.ErrorJson
	err = json.Unmarshal(resp, &httpError)
	if err != nil {
		return "", err, ""
	}
	fmt.Printf("Account create failed: %s", httpError.Message)

	return "", nil, httpError.Message
}
