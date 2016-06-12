package handlers

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/deployithq/deployit/drivers/interfaces"
	print_ "github.com/deployithq/deployit/drivers/print"
	"github.com/deployithq/deployit/utils"
	"net/http"
	"os"
	"path/filepath"
)

type ServiceCommand struct {
	ServiceName string
	Subcommand  string
	Host        struct {
		Name string
		URL  string
	}
	Paths struct {
		Root    string
		Storage string
	}
	Print interfaces.IPrint
}

func (c *ServiceCommand) Init(args []string) error {

	var err error
	var debug bool

	// Initializaing printing module
	c.Print = print_.Init()

	// Adding path module
	c.Paths.Root, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		c.Print.Error(err)
		return err
	}

	c.Paths.Storage = fmt.Sprintf("%s/.dit", c.Paths.Root)

	// Creating flags set
	cmdFlags := flag.NewFlagSet("service", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		c.Print.WhiteInfo(c.Help())
	}

	cmdFlags.BoolVar(&debug, "debug", false, "Enables debug mode")
	if Debug == false {
		if os.Getenv("DEPLOYIT_DEBUG") != "" {
			Debug = true
		}
	}

	cmdFlags.StringVar(&c.Host.URL, "host-url", "https://api.deployit.io", "URL of host, where daemon is running")
	if Debug == false {
		if os.Getenv("DEPLOYIT_HOST_URL") != "" {
			Host = os.Getenv("DEPLOYIT_HOST_URL")
		}
	}

	return nil

}

func (c *ServiceCommand) Run(args []string) int {

	err := c.Init(args)
	if err != nil {
		c.Print.Error(err)
		return 1
	}

	var res *http.Response

	if err != nil {
		c.Print.Error(err)
		return 1
	}

	// TODO Services switcher

	switch c.Subcommand {
	case "stop":
		ServiceStop(c)
	}

	if res.StatusCode != 200 {
		err = errors.New("Something went wrong")
		c.Print.Error(err)
		return 1
	}

	c.Print.Info("Finished!")

	return 0

}

func ServiceStop(c *ServiceCommand) (*http.Response, error) {

	c.Print.Infof("Stopping %s ...", c.ServiceName)

	res, err := utils.ExecuteHTTPRequest(fmt.Sprintf("%s/service/%s/stop", c.Host.URL, c.ServiceName), "POST", "application/json", new(bytes.Buffer))
	if err != nil {
		c.Print.Error(err)
		return res, err
	}

	return res, nil
}

func ServiceRemove(c *ServiceCommand, appInfo *AppInfo) (*http.Response, error) {

	c.Print.Infof("Removing %s ...", appInfo.Name)

	res, err := utils.ExecuteHTTPRequest(fmt.Sprintf("%s/service/%s", c.Host.URL, c.ServiceName), "DELETE", "application/json", new(bytes.Buffer))
	if err != nil {
		c.Print.Error(err)
		return res, err
	}

	return res, nil
}

func (c *ServiceCommand) Help() string {
	return ""
}

func (c *ServiceCommand) Synopsis() string {
	return ""
}

//func ServiceStart(c *cli.Context) error {
//
//	env := NewEnv()
//
//	if ServiceName == "" {
//		env.Log.Error("Unknown service")
//		return nil
//	}
//
//	color.Cyan("Starting %s ...", ServiceName)
//
//	res, err := http.Post(fmt.Sprintf("%s/service/%s/start", env.HostUrl, ServiceName), "application/json", new(bytes.Buffer))
//	if err != nil {
//		env.Log.Error(err)
//		return err
//	}
//
//	if res.StatusCode != 200 {
//		err = errors.ParseError(res)
//		env.Log.Error(err)
//		return err
//	}
//
//	response := struct {
//		Port     int64  `json:"port"`
//		Password string `json:"password"`
//	}{}
//
//	err = json.NewDecoder(res.Body).Decode(&response)
//	if err != nil {
//		env.Log.Error(err)
//		return err
//	}
//
//	color.Cyan("Your %s address: %s:%d", ServiceName, env.Host, response.Port)
//	color.Cyan("Your %s password: %s", ServiceName, response.Password)
//	color.Cyan("Finished!")
//
//	return nil
//}
//
//func ServiceRestart(c *cli.Context) error {
//
//	env := NewEnv()
//
//	if ServiceName == "" {
//		env.Log.Error("Unknown service")
//		return nil
//	}
//
//	color.Cyan("Restarting %s ...", ServiceName)
//
//	res, err := http.Post(fmt.Sprintf("%s/service/%s/restart", env.HostUrl, ServiceName), "application/json", new(bytes.Buffer))
//	if err != nil {
//		env.Log.Error(err)
//		return err
//	}
//
//	if res.StatusCode != 200 {
//		err = errors.ParseError(res)
//		env.Log.Error(err)
//		return err
//	}
//
//	response := struct {
//		Port     int64  `json:"port"`
//		Password string `json:"password"`
//	}{}
//
//	err = json.NewDecoder(res.Body).Decode(&response)
//	if err != nil {
//		env.Log.Error(err)
//		return err
//	}
//
//	color.Cyan("Your %s address: %s:%d", ServiceName, env.Host, response.Port)
//	color.Cyan("Your %s password: %s", ServiceName, response.Password)
//	color.Cyan("Finished!")
//
//	return nil
//}
//
//func ServiceDeploy(c *cli.Context) error {
//
//	// TODO Adapt for other services
//
//	env := NewEnv()
//
//	if ServiceName == "" {
//		env.Log.Error("Unknown service")
//		return nil
//	}
//
//	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/service/%s", env.HostUrl, ServiceName), new(bytes.Buffer))
//	if err != nil {
//		env.Log.Error(err)
//		return err
//	}
//
//	req.Header.Set("Content-Type", "application/json; charset=utf-8")
//
//	client := new(http.Client)
//	res, err := client.Do(req)
//	if err != nil {
//		env.Log.Error(err)
//		return err
//	}
//
//	if res.StatusCode != 200 {
//		err = errors.ParseError(res)
//		env.Log.Error(err)
//		return err
//	}
//
//	response := struct {
//		Port     int64  `json:"port"`
//		Password string `json:"password"`
//	}{}
//
//	err = json.NewDecoder(res.Body).Decode(&response)
//	if err != nil {
//		env.Log.Error(err)
//		return err
//	}
//
//	color.Cyan("Your %s address: %s:%d", ServiceName, env.Host, response.Port)
//	color.Cyan("Your %s password: %s", ServiceName, response.Password)
//
//	return nil
//}
