package utils

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
)

func ExecuteHTTPRequest(url string, method, contentType string, buf io.ReadWriter) (*http.Response, error) {

	res := new(http.Response)

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return res, err
	}

	req.Header.Set("Content-Type", contentType)

	client := new(http.Client)
	res, err = client.Do(req)
	if err != nil {
		return res, err
	}

	return res, nil

}

func StreamHttpResponse(res *http.Response) {

	reader := bufio.NewReader(res.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		fmt.Println(string(line))
	}

}
