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

package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"

	"github.com/lastbackend/lastbackend/internal/cli/command/options"
)

const defaultEnvPrefix = "LB"
const defaultConfigType = "yaml"
const defaultConfigName = "config"
const componentLB = "lb"

func New() *cobra.Command {

	global := pflag.CommandLine

	cleanFlagSet := pflag.NewFlagSet(componentLB, pflag.ContinueOnError)
	cleanFlagSet.SetNormalizeFunc(WordSepNormalizeFunc)

	var cmd = &cobra.Command{
		Use:   "lb",
		Short: "Last.Backend Open-source PaaS",
		Long:  `Open-source system for automating deployment, scaling, and management of containerized applications.`,
		// Because has special flag parsing requirements to enforce flag precedence rules,
		// so we do all our parsing manually in Run, below.
		// DisableFlagParsing=true provides the full set of flags passed
		// `args`  to Run, without Cobra's interference.
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {

			if err := cleanFlagSet.Parse(args); err != nil {
				cmd.Usage()
				fmt.Println(err)
				return
			}

			// check if there are non-flag arguments in the cmd line
			cmds := cleanFlagSet.Args()
			if len(cmds) > 0 {
				cmd.Usage()
				fmt.Println("unknown cmd: ", cmds[0])
				return
			}

			// short-circuit on help
			help, err := cleanFlagSet.GetBool("help")
			if err != nil {
				fmt.Println(`"help" flag is non-bool, programmer error, please correct`)
			}
			if help {
				cmd.Help()
				return
			}

			return
		},
	}

	global.BoolP("debug", "d", false, "Enable debug mode")

	options.AddGlobalFlags(cleanFlagSet)

	return cmd
}

// PrintFlags logs the flags in the flagset
func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		fmt.Println(fmt.Sprintf("FLAG: --%s=%q", flag.Name, flag.Value))
	})
}

// WordSepNormalizeFunc normalizes cli flags
func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}
