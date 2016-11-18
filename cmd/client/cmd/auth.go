package cmd

import (
)
import (
	"fmt"
	"github.com/howeyc/gopass"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"github.com/lastbackend/lastbackend/cmd/client/config"
)

type loginInfo struct {
	Login string `json:"login"`
	Password string `json:"password"`
}


func Auth() {
	fmt.Print("Login: ")
	var login string
	fmt.Scan(&login)

	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	password := string(pass)

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