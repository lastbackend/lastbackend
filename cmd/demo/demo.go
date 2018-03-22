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

package main

import (
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:  "app",
		Long: "QWEQWE",
		Run:  func(cmd *cobra.Command, args []string) {},
	}

	helpCmd := &cobra.Command{
		Use:   "help [command]",
		Short: "Help about any command",
		Long: `Help provides help for any command in the application.
	Simply type ` + cmd.Name() + ` help [path to command] for full details.`,

		Run: func(c *cobra.Command, args []string) {
			cmd, _, e := c.Root().Find(args)
			if cmd == nil || e != nil {
				c.Printf("Unknown help topic %#q\n", args)
				c.Root().Usage()
			} else {
				cmd.InitDefaultHelpFlag() // make possible 'help' flag to be shown
				cmd.Help()
			}
		},
	}

	helpCmd.Flags().BoolP("aflag", "a", true, "a")
	helpCmd.Flags().BoolP("lflag", "l", true, "l")
	cmd.SetHelpCommand(helpCmd)
	cmd.Execute()
}
