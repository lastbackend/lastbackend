package handlers

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/urfave/cli.v2"
	"net/http"
)

func AppStart(c *cli.Context) error {

	env := NewEnv()

	appInfo := new(AppInfo)
	err := appInfo.Read(env.Log, env.Path, env.Host)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	color.Cyan("Starting %s ...", appInfo.Name)

	res, err := http.Get(fmt.Sprintf("%s/%s/start", env.HostUrl, appInfo.UUID))
	if err != nil {
		env.Log.Error(err)
		return err
	}

	if res.StatusCode != 200 {
		err = errors.New("Something went wrong")
		env.Log.Error(err)
		return err
	}

	color.Cyan("Finished!")

	return nil

}

func AppStop(c *cli.Context) error {

	env := NewEnv()

	appInfo := new(AppInfo)
	err := appInfo.Read(env.Log, env.Path, env.Host)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	color.Cyan("Stopping %s ...", appInfo.Name)

	res, err := http.Get(fmt.Sprintf("%s/%s/stop", env.HostUrl, appInfo.UUID))
	if err != nil {
		env.Log.Error(err)
		return err
	}

	if res.StatusCode != 200 {
		err = errors.New("Something went wrong")
		env.Log.Error(err)
		return err
	}

	color.Cyan("Finished!")

	return nil

}
