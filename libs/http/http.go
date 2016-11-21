package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HTTP struct{}

func (h *HTTP) Post(url string, json []byte, header string, headerType string) ([]byte, int) {
	req := h.request("POST", url, json)
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

func (h *HTTP) Get(url string, json []byte, header string, headerType string) ([]byte, int) {
	req := h.request("GET", url, nil)
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

func (h *HTTP) Delete(url string, header string, headerType string) int {
	req := h.request("DELETE", url, nil)
	req.Header.Add(header, headerType)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	return resp.StatusCode
}

func (h *HTTP) request(method string, url string, json []byte) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(json))
	if err != nil {
		fmt.Println(err.Error())
	}
	return req
}
