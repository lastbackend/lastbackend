package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/jarcoal/httpmock"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	"github.com/lastbackend/lastbackend/cmd/client/context"
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
	var password string
	var login string

	if ctx == context.Mock() {
		login = "lavr"
		password = "12345678"

		httpmock.Activate()
		defer httpmock.Deactivate()

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
	} else {
		fmt.Print("Login: ")
		fmt.Scan(&login)

		fmt.Print("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		password = string(pass)
	}

	data := loginInfo{login, password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	req, err := http.NewRequest("POST", config.Get().AuthUserUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var token tokenInfo
	err = json.Unmarshal(respContent, &token)
}
