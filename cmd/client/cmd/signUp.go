package cmd

import (
	"bytes"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	"io/ioutil"
	"k8s.io/client-go/1.5/pkg/util/json"
	"net/http"
)

type newUserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenInfo struct {
	Token string `json:"token"`
}

func SignUp() {
	fmt.Print("Username: ")
	var username string
	fmt.Scan(&username)

	fmt.Print("Email: ")
	var email string
	fmt.Scan(&email)

	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	password := string(pass)

	data := newUserInfo{username, email, password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	req, err := http.NewRequest("POST", config.Get().CreateUserUrl, bytes.NewBuffer(jsonData))
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
	byteToken := []byte(token.Token)

	//TODO что за рандомный параметр 0544
	err = ioutil.WriteFile(config.Get().StoragePath, byteToken, 0544)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
