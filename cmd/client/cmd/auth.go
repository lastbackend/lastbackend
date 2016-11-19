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
	"net/http"
)

type loginInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type authToken struct {
	Token string `json:"token"`
}

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
		func(req *http.Request) (*http.Response, error) {

			article := authToken{
				Token: "token",
			}

			if err := json.NewDecoder(req.Body).Decode(&article); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}

			resp, err := httpmock.NewJsonResponse(200, article)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	return login, password
}
