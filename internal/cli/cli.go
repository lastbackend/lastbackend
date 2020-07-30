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

package cli

import (
	"fmt"
	"os"

	"github.com/lastbackend/lastbackend/internal/cli/command"
	"github.com/lastbackend/lastbackend/internal/cli/command/client"
	"github.com/lastbackend/lastbackend/internal/cli/command/daemon"
	"github.com/spf13/cobra"
)

type CLI struct {
	rootCmd *cobra.Command
}

func New() *CLI {
	c := new(CLI)
	rootCmd := command.New()
	rootCmd.AddCommand(command.VersionCmd)
	rootCmd.AddCommand(daemon.NewCommand())
	rootCmd.AddCommand(client.NewCommands()...)
	c.rootCmd = rootCmd
	return c
}

func (c *CLI) Execute() {
	if err := c.rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
