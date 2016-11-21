package mocks

import (
	"fmt"
	"io/ioutil"
	"bytes"
	"net/http"
	"github.com/jarcoal/httpmock"
)

type HTTPMock struct{}

func (h *HTTPMock) Post(url string, json []byte, header string, headerType string) ([]byte, int) {
	req := h.request("POST", "/project", json)
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

func (h *HTTPMock) Get(url string, json []byte, header string, headerType string) ([]byte, int) {
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
/*
func TestFetchArticles(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.mybiz.com/articles.json",
		httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Article"}]`))

	// do stuff that makes a request to articles.json
}
*/



func (h *HTTPMock) Delete(url string, header string, headerType string) int {
	httpmock.Activate()
	req := h.request("DELETE", url, nil)
	req.Header.Add(header, headerType)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	return resp.StatusCode
}

func (h *HTTPMock) request(method string, url string, json []byte) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(json))
	if err != nil {
		fmt.Println(err.Error())
	}
	return req
}

