package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)


func POST(url string, json []byte, header string, headerType string) ([]byte, int) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		fmt.Println(err.Error())
	}
	req.Header.Add(header, headerType)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	return respContent, resp.StatusCode
}

func Get(url string, json []byte, header string, headerType string) ([]byte, int) {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(json))
	if err != nil {
		fmt.Println(err.Error())
	}
	req.Header.Add(header, headerType)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	return respContent, resp.StatusCode
}
