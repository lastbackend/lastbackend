//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package service

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/cli/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"github.com/unloop/gopipe"
	"io"
	"strconv"
	"strings"
)

type mapInfo map[string]serviceInfo
type serviceInfo struct {
	Pod       string
	Container string
}

type Writer struct {
	io.Writer
}

func (Writer) Write(p []byte) (int, error) {
	return fmt.Print(string(p))
}

func LogsServiceCmd(name string) {

	var (
		choice string = "0"
	)

	service, namespace, err := Inspect(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	var m = make(mapInfo)
	var index int = 0

	fmt.Println("Contaner logs:\n")

	for _, pod := range service.Pods {
		for _, container := range pod.Containers {
			fmt.Printf("[%d] %s\n", index, container.ID)

			m[strconv.Itoa(index)] = serviceInfo{
				Pod:       pod.Meta.Name,
				Container: container.ID,
			}
		}
		index++
	}

	if len(m) > 1 {
		for {
			fmt.Print("\nEnter container number for watch log or ^C for Exit: ")
			fmt.Scan(&choice)
			choice = strings.ToLower(choice)

			if _, ok := m[choice]; ok {
				break
			}

			fmt.Println("Number not correct!")
		}
	}

	reader, err := Logs(namespace, service.Meta.Name, m[choice].Pod, m[choice].Container)
	if err != nil {
		fmt.Println(err)
		return
	}

	stream.New(Writer{}).Pipe(reader)
}

func Logs(namespace, name, pod, container string) (*io.ReadCloser, error) {

	var (
		err  error
		http = context.Get().GetHttpClient()
		er   = new(errors.Http)
	)

	_, res, err := http.
		GET("/namespace/" + namespace + "/service/" + name + "/logs?pod=" + pod + "&container=" + container).Do()
	if err != nil {
		return nil, err
	}

	if er.Code == 401 {
		return nil, errors.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, errors.New(er.Message)
	}

	return &res.Body, nil
}
