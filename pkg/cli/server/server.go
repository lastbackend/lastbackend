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

package server

import (
	"fmt"
	"github.com/lastbackend/lastbackend/internal/api"
	"github.com/lastbackend/lastbackend/internal/controller"
	"github.com/lastbackend/lastbackend/pkg/cli/server/options"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

const (
	componentServer = "server"
)

func NewServerCommand() *cobra.Command {

	global := pflag.CommandLine

	cleanFlagSet := pflag.NewFlagSet(componentServer, pflag.ContinueOnError)
	cleanFlagSet.SetNormalizeFunc(WordSepNormalizeFunc)
	serverFlags := options.NewServerFlags()

	var command = &cobra.Command{
		Use:                "server",
		Short:              "Run master component",
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
				fmt.Println("unknown command: %s", cmds[0])
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

			// TODO: configure viper
			//cmd.MarkFlagRequired("vault-token")
			//cmd.MarkFlagRequired("vault-endpoint")
			//
			//v.BindPFlag("token", fs.Lookup("access-token"))
			//v.BindPFlag("name", fs.Lookup("cluster-name"))
			//v.BindPFlag("description", fs.Lookup("cluster-description"))
			//v.BindPFlag("server.host", fs.Lookup("bind-address"))
			//v.BindPFlag("server.port", fs.Lookup("bind-port"))
			//v.BindPFlag("server.tls.cert", fs.Lookup("tls-cert-file"))
			//v.BindPFlag("server.tls.key", fs.Lookup("tls-private-key-file"))
			//v.BindPFlag("server.tls.ca", fs.Lookup("tls-ca-file"))
			//v.BindPFlag("vault.token", fs.Lookup("vault-token"))
			//v.BindPFlag("vault.endpoint", fs.Lookup("vault-endpoint"))
			//v.BindPFlag("domain.internal", fs.Lookup("domain-internal"))
			//v.BindPFlag("domain.external", fs.Lookup("domain-external"))
			//v.BindPFlag("storage.driver", fs.Lookup("storage"))
			//v.BindPFlag("storage.etcd.tls.cert", fs.Lookup("etcd-cert-file"))
			//v.BindPFlag("storage.etcd.tls.key", fs.Lookup("etcd-private-key-file"))
			//v.BindPFlag("storage.etcd.tls.ca", fs.Lookup("etcd-ca-file"))
			//v.BindPFlag("storage.etcd.endpoints", fs.Lookup("etcd-endpoints"))
			//v.BindPFlag("storage.etcd.prefix", fs.Lookup("etcd-prefix"))
			Run(viper.New())
		},
	}

	global.IntP("verbose", "v", 0, "Set log level from 0 to 7")

	serverFlags.AddFlags(cleanFlagSet)
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

func Run(v *viper.Viper) {

	var (
		done    = make(chan bool, 1)
		apps    = make(chan bool)
		wait    = 0
		daemons = []func(*viper.Viper) bool{
			api.Daemon,
			controller.Daemon,
		}
	)

	for _, app := range daemons {
		go func(app func(*viper.Viper) bool, v *viper.Viper) {
			wait++
			apps <- app(v)
		}(app, v)
	}

	go func() {
		for {
			select {
			case <-apps:
				wait--
				if wait == 0 {
					done <- true
					return
				}
			}
		}
	}()

	<-done
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
