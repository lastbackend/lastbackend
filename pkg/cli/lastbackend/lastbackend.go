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
	"github.com/lastbackend/lastbackend/pkg/cli/agent"
	"os"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/cli/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultEnvPrefix = "LB"
const defaultConfigType = "yaml"
const defaultConfigName = "config"

var cfgFile string

func NewLbCommand() *cobra.Command {


	var command = &cobra.Command{
		Use:   "lb",
		Short: "Last.Backend Open-source API",
		Long:  `Open-source system for automating deployment, scaling, and management of containerized applications.`,
		Run: func(cmd *cobra.Command, args []string) {
			server.Run(viper.New())
			agent.Run(viper.New())
		},
	}

	command.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Path for the configuration file")
	command.PersistentFlags().IntP("verbose", "v", 0, "Set log level from 0 to 7")

	viper.BindPFlag("service.cidr", command.PersistentFlags().Lookup("services-cidr"))
	viper.BindPFlag("verbose", command.PersistentFlags().Lookup("verbose"))

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
