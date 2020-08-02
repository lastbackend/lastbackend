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
	"context"
	"fmt"

	"github.com/howeyc/gopass"
	"github.com/lastbackend/lastbackend/internal/cli/models"
	"github.com/lastbackend/lastbackend/internal/cli/service"
	"github.com/lastbackend/lastbackend/pkg/client/genesis/http/v1/request"
	"github.com/lastbackend/lastbackend/tools/logger"
	"github.com/spf13/cobra"
)

const logInExample = `
  # Log in to a Last.Backend 
  lb login
  Login: username
  Password: ******
`

const logOutExample = `
  # Log out from a Last.Backend 
  lb logout"
`

func (c *command) NewSessionLogInCmd(sessionService *service.SessionService) *cobra.Command {
	log := logger.WithContext(context.Background())

	return &cobra.Command{
		Use:     "login",
		Short:   "Log in to a Last.Backend",
		Example: logInExample,
		Run: func(cmd *cobra.Command, args []string) {

			var (
				login    string
				password string
			)

			fmt.Print("Login: ")
			if _, err := fmt.Scan(&login); err != nil {
				log.Error(err.Error())
				return
			}

			fmt.Print("Password: ")
			pass, err := gopass.GetPasswd()
			if err != nil {
				log.Error(err.Error())
				return
			}

			password = string(pass)
			fmt.Print("\r\n")

			opts := &request.AccountLoginOptions{
				Login:    login,
				Password: password,
			}

			res, err := c.client.genesis.V1().Account().Login(context.Background(), opts)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			session := new(models.Session)
			session.Token = res.Token

			if err := sessionService.Set(session); err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Println("Authorization successful!")
			return
		},
	}
}

func (c *command) NewSessionLogOutCmd(sessionService *service.SessionService) *cobra.Command {
	return &cobra.Command{
		Use:     "logout",
		Short:   "Log out from a Last.Backend",
		Example: logOutExample,
		Run: func(cmd *cobra.Command, args []string) {
			if err := sessionService.Del(); err != nil {
				fmt.Println(err)
				return
			}
		},
	}
}
