//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lastbackend/lastbackend/tools/logger"
	"io"
	"os"
	"strings"

	"github.com/lastbackend/lastbackend/internal/cli/views"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/spf13/cobra"
)

const jobListExample = `
  # Get all jobs for 'ns-demo' namespace  
  lb job ls ns-demo
`

const jobInspectExample = `
  # Get information for 'redis' job in 'ns-demo' namespace
  lb job inspect ns-demo redis
`

const jobRunExample = `
  # Get information for 'redis' job in 'ns-demo' namespace
  lb job run ns/cron
`

const jobRemoveExample = `
  # Remove 'redis' job in 'ns-demo' namespace
  lb job remove ns-demo redis
`

const jobLogsExample = `
  # Get 'redis' job logs for 'ns-demo' namespace
  lb job logs [NAMESPACE]/[NAME] -t [task-id]
`

func (c *command) NewJobCmd() *cobra.Command {

	log := logger.WithContext(context.Background())

	cmd := &cobra.Command{
		Use:   "job",
		Short: "Manage your job",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Error(err.Error())
				return
			}
		},
	}

	cmd.AddCommand(c.jobListCmd())
	cmd.AddCommand(c.jobInspectCmd())
	cmd.AddCommand(c.jobRunCmd())
	cmd.AddCommand(c.jobRemoveCmd())
	cmd.AddCommand(c.jobLogsCmd())

	return cmd
}

func (c *command) jobListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls [NAMESPACE]",
		Short:   "Display the jobs list",
		Example: jobListExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]

			response, err := c.client.cluster.V1().Namespace(namespace).Job().List(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			if response == nil || len(*response) == 0 {
				fmt.Println("no jobs available")
				return
			}

			list := views.FromApiJobListView(response)
			list.Print()
		},
	}

	return cmd
}

func (c *command) jobInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inspect [NAMESPACE]/[NAME]",
		Short:   "Service info by name",
		Example: jobInspectExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			namespace, name, err := jobParseSelfLink(args[0])
			checkError(err)

			t, err := cmd.Flags().GetString("task")
			if err != nil {
				_ = fmt.Errorf("can not be parse task option: %s", t)
				return
			}

			if t == "" {
				job, err := c.client.cluster.V1().Namespace(namespace).Job(name).Get(context.Background())
				if err != nil {
					fmt.Println(err)
					return
				}

				ss := views.FromApiJobView(job)
				ss.Print()
				return
			}

			task, err := c.client.cluster.V1().Namespace(namespace).Job(name).Task(t).Get(context.Background())

			tw := views.FromApiTaskView(task)
			tw.Print()

		},
	}

	cmd.Flags().StringP("task", "t", "", "inspect particular task")

	return cmd
}

func (c *command) jobRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run [NAMESPACE]/[NAME]",
		Short:   "Run job info by name",
		Example: jobRunExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			namespace, name, err := jobParseSelfLink(args[0])
			checkError(err)

			_, err = c.client.cluster.V1().Namespace(namespace).Job(name).Run(context.Background(), nil)
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}

	return cmd
}

func (c *command) jobRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove [NAMESPACE] [NAME]",
		Short:   "Remove job by name",
		Example: jobRemoveExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			namespace := args[0]
			name := args[1]

			opts := &request.JobRemoveOptions{Force: false}

			if err := opts.Validate(); err != nil {
				fmt.Println(err.Err())
				return
			}

			if err := c.client.cluster.V1().Namespace(namespace).Job(name).Remove(context.Background(), opts); err != nil {
				_ = fmt.Errorf("job remove err: %s", err.Error())
				return
			}

			fmt.Println(fmt.Sprintf("Job `%s` is successfully removed", name))
		},
	}

	return cmd
}

func (c *command) jobLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logs [NAMESPACE]/[NAME]",
		Short:   "Get job logs",
		Example: jobLogsExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			opts := new(request.JobLogsOptions)

			namespace, name, err := jobParseSelfLink(args[0])
			checkError(err)

			task, err := cmd.Flags().GetString("task")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			opts.Tail, err = cmd.Flags().GetInt("tail")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			opts.Follow, err = cmd.Flags().GetBool("follow")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if task != models.EmptyString {
				opts.Task = task
			}

			reader, _, err := c.client.cluster.V1().Namespace(namespace).Job(name).Logs(context.Background(), opts)
			if err != nil {
				fmt.Println(err)
				return
			}

			dec := json.NewDecoder(reader)
			for {
				var doc models.LogMessage

				err := dec.Decode(&doc)
				if err == io.EOF {
					// all done
					break
				}
				if err != nil {
					fmt.Errorf(err.Error())
					os.Exit(1)
				}

				if doc.ContainerType == models.ContainerTypeRuntimeTask {
					fmt.Println(doc.Data)
				}
			}
		},
	}

	cmd.Flags().IntP("tail", "t", 0, "tail last n lines")
	cmd.Flags().BoolP("follow", "f", false, "follow logs")
	cmd.Flags().String("task", "", "read logs for particular task")

	return cmd
}

func jobParseSelfLink(selflink string) (string, string, error) {
	match := strings.Split(selflink, "/")

	var (
		namespace, name string
	)

	switch len(match) {
	case 2:
		namespace = match[0]
		name = match[1]
	case 1:
		fmt.Println("Use default namespace:", models.DEFAULT_NAMESPACE)
		namespace = models.DEFAULT_NAMESPACE
		name = match[0]
	default:
		return "", "", errors.New("invalid service name provided")
	}

	return namespace, name, nil
}
