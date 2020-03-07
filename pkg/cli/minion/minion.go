//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package minion

import (
	"fmt"
	"github.com/lastbackend/lastbackend/internal/minion"
	"github.com/lastbackend/lastbackend/pkg/cli/minion/options"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const defaultEnvPrefix = "LB"
const defaultConfigType = "yaml"
const defaultConfigName = "config"
const componentAgent = "agent"

var cfgFile string

func NewMinionCommand() *cobra.Command {

	initConfig()

	global := pflag.CommandLine

	cleanFlagSet := pflag.NewFlagSet(componentAgent, pflag.ContinueOnError)
	cleanFlagSet.SetNormalizeFunc(WordSepNormalizeFunc)
	agentFlags := options.NewMinionFlags()

	var command = &cobra.Command{
		Use:                "agent",
		Short:              "Run node agent",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {

			if err := cleanFlagSet.Parse(args); err != nil {
				cmd.Usage()
				fmt.Println(err)
				return
			}

			// check if there are non-flag arguments in the command line
			cmds := cleanFlagSet.Args()
			if len(cmds) > 0 {
				cmd.Usage()
				fmt.Println("unknown command: ", cmds[0])
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

			PrintFlags(cleanFlagSet)

			Run(viper.New())
		},
	}

	global.StringVarP(&cfgFile, "config", "c", "", "Path for the configuration file")
	global.IntP("verbose", "v", 0, "Set log level from 0 to 7")

	agentFlags.AddFlags(cleanFlagSet)
	options.AddGlobalFlags(cleanFlagSet)

	cleanFlagSet.BoolP("help", "h", false, fmt.Sprintf("help for %s", command.Name()))

	const usageFmt = "Usage:\n  %s\n\nFlags:\n%s"
	command.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine(), cleanFlagSet.FlagUsagesWrapped(2))
		return nil
	})

	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine(), cleanFlagSet.FlagUsagesWrapped(2))
	})

	return command
}

func initConfig() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix(defaultEnvPrefix)
	viper.SetConfigType(defaultConfigType)
	viper.SetConfigFile(viper.GetString(defaultConfigName))

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Can't read config:", err)
			os.Exit(1)
		}
	}
}

func Run(v *viper.Viper) {
	minion.Daemon(v)
}

// PrintFlags logs the flags in the flagset
func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		fmt.Println(fmt.Sprintf("FLAG: --%s=%q", flag.Name, flag.Value))
	})
}

func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}
