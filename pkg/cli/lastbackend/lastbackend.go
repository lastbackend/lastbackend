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

package lastbackend

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/cli/lastbackend/options"
	"os"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/cli/master"
	mo "github.com/lastbackend/lastbackend/pkg/cli/master/options"
	"github.com/lastbackend/lastbackend/pkg/cli/minion"
	no "github.com/lastbackend/lastbackend/pkg/cli/minion/options"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const defaultEnvPrefix = "LB"
const defaultConfigType = "yaml"
const defaultConfigName = "config"
const componentLB = "lastbackend"

var cfgFile string

// NewLBCmd entrypoint for CLI launcher
func NewRootCommand() *cobra.Command {

	initConfig()

	global := pflag.CommandLine

	cleanFlagSet := pflag.NewFlagSet(componentLB, pflag.ContinueOnError)
	cleanFlagSet.SetNormalizeFunc(WordSepNormalizeFunc)
	masterFlags := mo.NewMasterFlags()
	minionFlags := no.NewMinionFlags()

	var command = &cobra.Command{
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

			masterViper := viper.New()
			if len(cfgFile) == 0 {
				masterViper = masterFlags.LoadViper(viper.New())
			}

			minionViper := viper.New()
			if len(cfgFile) == 0 {
				minionViper = minionFlags.LoadViper(viper.New())
			}

			master.Run(masterViper)
			minion.Run(minionViper)
		},
	}

	global.StringVarP(&cfgFile, "config", "c", "", "Path for the configuration file")
	global.IntP("verbose", "v", 0, "Set log level from 0 to 7")

	// keep cleanFlagSet separate, so Cobra doesn't pollute it with the global flags
	masterFlags.AddFlags(cleanFlagSet)
	minionFlags.AddFlags(cleanFlagSet)
	options.AddGlobalFlags(cleanFlagSet)

	cleanFlagSet.BoolP("help", "h", false, fmt.Sprintf("help for %s", command.Name()))

	// this necessary, because Cobra's default UsageFunc and HelpFunc pollute the flagset with global flags
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
