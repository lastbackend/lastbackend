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

package agent

import (
	"fmt"
	"github.com/lastbackend/lastbackend/internal/agent"
	"github.com/lastbackend/lastbackend/pkg/cli/agent/options"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

const (
	componentAgent = "agent"
)

func NewServerCommand() *cobra.Command {

	cleanFlagSet := pflag.NewFlagSet(componentAgent, pflag.ContinueOnError)
	cleanFlagSet.SetNormalizeFunc(WordSepNormalizeFunc)
	serverFlags := options.NewServerFlags()

	var command = &cobra.Command{
		Use:                "agent",
		Short:              "Run node agent",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {

			if err := cleanFlagSet.Parse(args); err != nil {
				cmd.Usage()
				fmt.Println(err)
			}

			// check if there are non-flag arguments in the command line
			cmds := cleanFlagSet.Args()
			if len(cmds) > 0 {
				cmd.Usage()
				fmt.Println("unknown command: %s", cmds[0])
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
			//v.BindPFlag("token", fs.Lookup("access-token"))
			//v.BindPFlag("workdir", fs.Lookup("workdir"))
			//v.BindPFlag("manifest.dir", fs.Lookup("manifest-path"))
			//v.BindPFlag("network.interface", fs.Lookup("bind-interface"))
			//v.BindPFlag("network.cpi.type", fs.Lookup("network-proxy"))
			//v.BindPFlag("network.cpi.interface.internal", fs.Lookup("network-proxy-iface-internal"))
			//v.BindPFlag("network.cpi.interface.external", fs.Lookup("network-proxy-iface-external"))
			//v.BindPFlag("network.cni.type", fs.Lookup("network-driver"))
			//v.BindPFlag("network.cni.interface.external", fs.Lookup("network-driver-iface-external"))
			//v.BindPFlag("network.cni.interface.internal", fs.Lookup("network-driver-iface-internal"))
			//v.BindPFlag("container.cri.type", fs.Lookup("container-runtime"))
			//v.BindPFlag("container.cri.docker.version", fs.Lookup("container-runtime-docker-version"))
			//v.BindPFlag("container.cri.docker.host", fs.Lookup("container-runtime-docker-host"))
			//v.BindPFlag("container.cri.docker.tls.verify", fs.Lookup("container-runtime-docker-tls-verify"))
			//v.BindPFlag("container.cri.docker.tls.ca_file", fs.Lookup("container-runtime-docker-tls-ca"))
			//v.BindPFlag("container.cri.docker.tls.cert_file", fs.Lookup("container-runtime-docker-tls-cert"))
			//v.BindPFlag("container.cri.docker.tls.key_file", fs.Lookup("container-runtime-docker-tls-key"))
			//v.BindPFlag("container.csi.dir.root", fs.Lookup("container-storage-root"))
			//v.BindPFlag("container.iri.type", fs.Lookup("container-image-runtime"))
			//v.BindPFlag("container.iri.docker.version", fs.Lookup("container-image-runtime-docker-version"))
			//v.BindPFlag("container.iri.docker.host", fs.Lookup("container-image-runtime-docker-host"))
			//v.BindPFlag("container.iri.docker.tls.verify", fs.Lookup("container-image-runtime-docker-tls-verify"))
			//v.BindPFlag("container.iri.docker.tls.ca_file", fs.Lookup("container-image-runtime-docker-tls-ca"))
			//v.BindPFlag("container.iri.docker.tls.cert_file", fs.Lookup("container-image-runtime-docker-tls-cert"))
			//v.BindPFlag("container.iri.docker.tls.key_file", fs.Lookup("container-image-runtime-docker-tls-key"))
			//v.BindPFlag("container.extra_hosts", fs.Lookup("container-extra-hosts"))
			//v.BindPFlag("server.host", fs.Lookup("bind-address"))
			//v.BindPFlag("server.port", fs.Lookup("bind-port"))
			//v.BindPFlag("server.tls.verify", fs.Lookup("tls-verify"))
			//v.BindPFlag("server.tls.cert", fs.Lookup("tls-cert-file"))
			//v.BindPFlag("server.tls.key", fs.Lookup("tls-private-key-file"))
			//v.BindPFlag("server.tls.ca", fs.Lookup("tls-ca-file"))
			//v.BindPFlag("api.uri", fs.Lookup("api-uri"))
			//v.BindPFlag("api.tls.verify", fs.Lookup("api-tls-verify"))
			//v.BindPFlag("api.tls.cert", fs.Lookup("api-tls-cert-file"))
			//v.BindPFlag("api.tls.key", fs.Lookup("api-tls-private-key-file"))
			//v.BindPFlag("api.tls.ca", fs.Lookup("api-tls-ca-file"))

			Run(viper.New())
		},
	}

	serverFlags.AddFlags(cleanFlagSet)

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
	agent.Daemon(v)
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
