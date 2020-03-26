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

package lastbackend

import (
	"fmt"
	"os"
	"strings"

	"github.com/lastbackend/lastbackend/internal/master"
	"github.com/lastbackend/lastbackend/internal/minion"
	"github.com/lastbackend/lastbackend/pkg/cli/lastbackend/options"
	mo "github.com/lastbackend/lastbackend/pkg/cli/master/options"
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

	var agent = false
	var noSchedule = false

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
			masterViper.AutomaticEnv()
			masterViper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			masterViper.SetEnvPrefix(defaultEnvPrefix)
			masterViper.SetConfigType(defaultConfigType)
			masterViper.SetConfigFile(masterViper.GetString(defaultConfigName))

			// Use config file from the flag.
			if cfgFile != "" {
				masterViper.SetConfigFile(cfgFile)
				if err := masterViper.ReadInConfig(); err != nil {
					fmt.Println(err)
					return
				}
			}

			minionViper := viper.New()
			minionViper.AutomaticEnv()
			minionViper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			minionViper.SetEnvPrefix(defaultEnvPrefix)
			minionViper.SetConfigType(defaultConfigType)
			minionViper.SetConfigFile(minionViper.GetString(defaultConfigName))

			// Use config file from the flag.
			if cfgFile != "" {
				minionViper.SetConfigFile(cfgFile)
				if err := minionViper.ReadInConfig(); err != nil {
					fmt.Println(err)
					return
				}
			}

			Run(masterViper, minionViper, &RunOptions{DisableMaster: agent, DisableMinion: noSchedule})
		},
	}

	global.StringVarP(&cfgFile, "config", "c", "", "Path for the configuration file")
	global.IntP("verbose", "v", 0, "Set log level from 0 to 7")
	global.BoolVarP(&agent, "agent", "", false, "Only agent mode")
	global.BoolVarP(&noSchedule, "no-schedule", "", false, "Disable schedule mode")

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

type RunOptions struct {
	DisableMaster bool
	DisableMinion bool
}

func Run(masterViper *viper.Viper, minionViper *viper.Viper, opts *RunOptions) {

	var (
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
		err  error
	)

	if opts == nil {
		opts = &RunOptions{false, false}
	}

	if opts.DisableMinion && opts.DisableMaster {
		fmt.Println("\n#################################")
		fmt.Println("### All services was disabled ###")
		fmt.Println("#################################\n")
		return
	}

	var masterApp *master.App
	var minionApp *minion.App

	if !opts.DisableMaster {
		masterApp, err = master.New(masterViper)
		if err != nil {
			panic(fmt.Sprintf("Create master application err: %v", err))
		}

		if err := masterApp.Run(); err != nil {
			panic(fmt.Sprintf("Run master application err: %v", err))
		}
	}

	if !opts.DisableMinion {
		minionApp, err = minion.New(minionViper)
		if err != nil {
			panic(fmt.Sprintf("Create minion application err: %v", err))
		}

		if err := minionApp.Run(); err != nil {
			panic(fmt.Sprintf("Run minion application err: %v", err))
		}
	}

	go func() {
		for {
			select {
			case <-sigs:
				if !opts.DisableMaster {
					masterApp.Stop()
				}
				if !opts.DisableMaster {
					minionApp.Stop()
				}
				done <- true
				return
			}
		}
	}()

	<-done
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
