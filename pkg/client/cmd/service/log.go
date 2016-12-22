package service

import (
	"errors"
	"fmt"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/pkg/client/context"
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
		ctx    = context.Get()
		choice string
	)

	service, projectName, err := Inspect(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	var m = make(mapInfo)
	var index int = 0

	fmt.Println("Contaner list:\n")

	for _, pod := range service.Detail.PodList.Pods {
		for _, container := range pod.ContainerList.Containers {
			fmt.Printf("[%d] %s\n", index, container.Name)

			m[strconv.Itoa(index)] = serviceInfo{
				Pod:       pod.ObjectMeta.Name,
				Container: container.Name,
			}
		}
		index++
	}

	for {
		fmt.Print("\nEnter container number for watch log or ^C for Exit: ")
		fmt.Scan(&choice)
		choice = strings.ToLower(choice)

		if _, ok := m[choice]; ok {
			break
		}

		ctx.Log.Error("Number not correct!")
	}

	reader, err := Logs(projectName, service.Name, m[choice].Pod, m[choice].Container)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	stream.New(Writer{}).Pipe(reader)
}

func Logs(project, name, pod, container string) (*io.ReadCloser, error) {

	var (
		err error
		ctx = context.Get()
		er  = new(e.Http)
	)

	_, res, err := ctx.HTTP.
		GET("/project/"+project+"/service/"+name+"/logs?pod="+pod+"&container="+container).
		AddHeader("Authorization", "Bearer "+ctx.Token).Do()
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if er.Code == 401 {
		return nil, errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	}

	if er.Code != 0 {
		return nil, errors.New(e.Message(er.Status))
	}

	return &res.Body, nil
}
