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
// patents in process, and are protected by trade secretCmd or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package cmd

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/cli/envs"
	"github.com/spf13/cobra"
	"github.com/unloop/gopipe"
	"io"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"strconv"
	"strings"
)

func init() {
	serviceCmd.AddCommand(serviceLogsCmd)
}

const serviceLogsExample = `
  # Get 'redis' service logs for 'ns-demo' namespace  
  lb service logs ns-demo redis
`

type Writer struct {
	io.Writer
}

func (Writer) Write(p []byte) (int, error) {
	return fmt.Print(string(p))
}

type mapInfo map[string]serviceInfo
type serviceInfo struct {
	Deployment string
	Pod        string
	Container  string
}

var serviceLogsCmd = &cobra.Command{
	Use:     "logs [NAMESPACE] [NAME]",
	Short:   "Get service logs",
	Example: serviceLogsExample,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			choice = "0"
			m      = make(mapInfo)
			index  = 0
		)

		namespace := args[0]
		name := args[1]

		cli := envs.Get().GetClient()
		response, err := cli.V1().Namespace(namespace).Service(name).Get(envs.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, deployment := range response.Deployments {
			for _, pod := range deployment.Pods {
				for _, container := range pod.Spec.Template.Containers {
					fmt.Printf("[%d] %s\n", index, container.Image.Name)
					m[strconv.Itoa(index)] = serviceInfo{
						Deployment: deployment.Meta.Name,
						Pod:        pod.Meta.Name,
						Container:  container.Image.Name,
					}
				}
				index++
			}
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

		opts := new(request.ServiceLogsOptions)
		opts.Deployment = m[choice].Deployment
		opts.Pod = m[choice].Pod
		opts.Container = m[choice].Container

		reader, err := cli.V1().Namespace(namespace).Service(name).Logs(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Service logs:\n")

		stream.New(Writer{}).Pipe(&reader)
	},
}
