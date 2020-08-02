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
	"fmt"
	"os"
	"path"

	"github.com/lastbackend/lastbackend/internal/pkg/storage"

	"github.com/lastbackend/lastbackend/internal/cli/service"
	"github.com/lastbackend/lastbackend/internal/util/filesystem"
	"github.com/lastbackend/lastbackend/pkg/client/cluster"
	"github.com/lastbackend/lastbackend/pkg/client/genesis"
	"github.com/spf13/cobra"
)

type command struct {
	client struct {
		cluster cluster.IClient
		genesis genesis.IClient
	}
}

func NewCommands() []*cobra.Command {
	cmd := new(command)

	if err := filesystem.MkDir(path.Join(filesystem.HomeDir(), "lastbackend"), 0755); err != nil {
		fmt.Printf("cannot create workdir: %v", err)
		os.Exit(1)
	}

	stg, err := storage.Get(storage.BboltDriver, storage.BboltConfig{DbDir: path.Join(filesystem.HomeDir(), "lastbackend"), DbName: ".lastbackend"})
	if err != nil {
		fmt.Printf("cannot initialize storage: %v", err)
		os.Exit(1)
	}

	sessionService := service.NewSessionService(stg)
	clusterService := service.NewClusterService(stg)

	return []*cobra.Command{
		cmd.NewSessionLogInCmd(sessionService),
		cmd.NewSessionLogOutCmd(sessionService),
		cmd.NewClusterCmd(clusterService),
		cmd.NewNodeCmd(),
		cmd.NewNamespaceCmd(),
		cmd.NewServiceCmd(),
		cmd.NewSecretCmd(),
		cmd.NewConfigCmd(),
		cmd.NewRouteCmd(),
		cmd.NewJobCmd(),
		cmd.NewVolumeCmd(),
		cmd.NewIngressCmd(),
		cmd.NewDiscoveryCmd(),
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func checkFlagParseError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
