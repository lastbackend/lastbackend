package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/deployithq/deployit/errors"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"net/http"
	"strconv"
)

func ServiceStart(c *cli.Context) error {

	env := NewEnv()

	if ServiceName == "" {
		env.Log.Error("Unknown service")
		return nil
	}

	color.Cyan("Starting %s ...", ServiceName)

	res, err := http.Post(fmt.Sprintf("%s/service/%s/start", env.HostUrl, ServiceName), "application/json", new(bytes.Buffer))
	if err != nil {
		env.Log.Error(err)
		return err
	}

	if res.StatusCode != 200 {
		err = errors.ParseError(res)
		env.Log.Error(err)
		return err
	}

	color.Cyan("Finished!")

	return nil
}

func ServiceStop(c *cli.Context) error {
	env := NewEnv()

	if ServiceName == "" {
		env.Log.Error("Unknown service")
		return nil
	}

	color.Cyan("Stopping %s ...", ServiceName)

	res, err := http.Post(fmt.Sprintf("%s/service/%s/stop", env.HostUrl, ServiceName), "application/json", new(bytes.Buffer))
	if err != nil {
		env.Log.Error(err)
		return err
	}

	if res.StatusCode != 200 {
		err = errors.ParseError(res)
		env.Log.Error(err)
		return err
	}

	color.Cyan("Finished!")

	return nil
}

func ServiceRestart(c *cli.Context) error {

	env := NewEnv()

	if ServiceName == "" {
		env.Log.Error("Unknown service")
		return nil
	}

	color.Cyan("Restarting %s ...", ServiceName)

	res, err := http.Post(fmt.Sprintf("%s/service/%s/restart", env.HostUrl, ServiceName), "application/json", new(bytes.Buffer))
	if err != nil {
		env.Log.Error(err)
		return err
	}

	if res.StatusCode != 200 {
		err = errors.ParseError(res)
		env.Log.Error(err)
		return err
	}

	color.Cyan("Finished!")

	return nil
}

func ServiceDeploy(c *cli.Context) error {

	// TODO Adapt for other services

	env := NewEnv()

	if ServiceName == "" {
		env.Log.Error("Unknown service")
		return nil
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/service/%s", env.HostUrl, ServiceName), new(bytes.Buffer))
	if err != nil {
		env.Log.Error(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	if res.StatusCode != 200 {
		err = errors.ParseError(res)
		env.Log.Error(err)
		return err
	}

	response := struct {
		Port     int64  `json:"port"`
		Password string `json:"password"`
	}{}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	color.Cyan("Your %s adress: %s:%s", ServiceName, env.HostUrl, strconv.FormatInt(response.Port, 10))
	color.Cyan("Your %s password: %s", ServiceName, response.Password)

	return nil
}

func ServiceRemove(c *cli.Context) error {

	env := NewEnv()

	if ServiceName == "" {
		env.Log.Error("Unknown service")
		return nil
	}

	color.Cyan("Removing %s ...", ServiceName)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/service/%s", env.HostUrl, ServiceName), new(bytes.Buffer))
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
		err = errors.ParseError(res)
		env.Log.Error(err)
		return err
	}

	color.Cyan("Finished!")

	return nil
}
