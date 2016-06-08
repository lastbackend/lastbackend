package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/urfave/cli.v2"
	"net/http"
	"os"
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

	res, err := http.Get(fmt.Sprintf("%s/app/%s/start", env.HostUrl, appInfo.UUID))
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

func AppRestart(c *cli.Context) error {

	env := NewEnv()

	appInfo := new(AppInfo)
	err := appInfo.Read(env.Log, env.Path, env.Host)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	color.Cyan("Restarting %s ...", appInfo.Name)

	res, err := http.Get(fmt.Sprintf("%s/app/%s/restart", env.HostUrl, appInfo.UUID))
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

	res, err := http.Get(fmt.Sprintf("%s/app/%s/stop", env.HostUrl, appInfo.UUID))
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

func AppRemove(c *cli.Context) error {

	env := NewEnv()

	appInfo := new(AppInfo)
	err := appInfo.Read(env.Log, env.Path, env.Host)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	color.Cyan("Removing %s ...", appInfo.Name)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/app/%s", env.HostUrl, appInfo.UUID), new(bytes.Buffer))
	if err != nil {
		env.Log.Error(err)
		return err
	}

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	if res.StatusCode != 200 {
		err = errors.New("Something went wrong")
		env.Log.Error(err)
		return err
	}

	os.Remove(fmt.Sprintf("%s/%s_map", env.StoragePath, env.Host))
	os.Remove(fmt.Sprintf("%s/%s.yaml", env.StoragePath, env.Host))

	// TODO Remove files connected with this host

	color.Cyan("Finished!")

	return nil

}
