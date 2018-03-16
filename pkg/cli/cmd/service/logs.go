//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	n "github.com/lastbackend/lastbackend/pkg/cli/cmd/namespace"
	e "github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

type mapInfo map[string]serviceInfo
type serviceInfo struct {
	Pod       string
	Container string
}

func LogsServiceCmd(name string) {

	var (
		choice = "0"
	)

	service, err := Inspect(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	var m = make(mapInfo)
	var index = 0

	fmt.Println("Service: ", service.Meta.Name)
	fmt.Println("Container logs: ")

	for _, dep := range service.Deployments {
		for _, pod := range dep.Pods {
			for _, con := range *pod.Status.Containers {

				fmt.Printf("[%d] %s\n", index, con.ID)

				m[strconv.Itoa(index)] = serviceInfo{
					Pod:       pod.ID,
					Container: con.ID,
				}
				index++
			}

		}
	}

	if len(m) > 1 {
		for {
			fmt.Print("\nEnter container number to watch the log or do ^C to Exit: ")
			fmt.Scan(&choice)
			choice = strings.ToLower(choice)

			if _, ok := m[choice]; ok {
				break
			}

			fmt.Println("Number isn't correct!")
		}
	}

	err = Logs(name, m[choice].Pod, m[choice].Container)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func Logs(name, pod, container string) error {

	var (
		err error
		URL = "wss://wss.lstbknd.net"
	)

	var dialer *websocket.Dialer

	ns, err := n.Current()
	if err != nil {
		return e.UnknownMessage
	}

	if ns.Meta == nil {
		return e.New("Workspace didn't select")
	}

	url := fmt.Sprintf("%s/namespace/%s/service/%s/%s/%s/logs", URL, ns.Meta.Name, name, pod, container)

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return e.UnknownMessage
	}

	for {

		var dat map[string]interface{}

		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return err
		}

		if err := json.Unmarshal(message, &dat); err != nil {
			return e.UnknownMessage
		}

		fmt.Printf("%s\n", dat["data"].(string))
	}

	return nil
}
