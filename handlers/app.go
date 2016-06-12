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

type AppCommand struct {
	Subcommand string
	Host       struct {
		Name string
		URL  string
	}
	Paths struct {
		Root    string
		Storage string
	}
	Print interfaces.IPrint
}

func Init(c *AppCommand, args []string) error {

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
	cmdFlags := flag.NewFlagSet("app", flag.ContinueOnError)
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

func (c *AppCommand) Run(args []string) int {

	err := Init(c, args)
	if err != nil {
		c.Print.Error(err)
		return 1
	}

	appInfo := new(AppInfo)
	err = appInfo.Read(c.Paths.Root, c.Host.Name)
	if err != nil {
		c.Print.Error(err)
		return 1
	}

	if appInfo.Name == "" {
		err = errors.New("App not found")
		c.Print.Error(err)
		return 1
	}

	var res *http.Response

	switch c.Subcommand {
	case "start":
		res, err = AppStart(c, appInfo)
	case "restart":
		res, err = AppRestart(c, appInfo)
	case "stop":
		res, err = AppStop(c, appInfo)
	case "remove":
		res, err = AppRemove(c, appInfo)
	}

	if err != nil {
		c.Print.Error(err)
		return 1
	}

	if res.StatusCode != 200 {
		err = errors.New("Something went wrong")
		c.Print.Error(err)
		return 1
	}

	c.Print.Info("Finished!")

	return 0

}

func AppStart(c *AppCommand, appInfo *AppInfo) (*http.Response, error) {

	c.Print.Infof("Starting %s ...", appInfo.Name)

	res, err := utils.ExecuteHTTPRequest(fmt.Sprintf("%s/app/%s/start", c.Host.URL, appInfo.Name), "POST", "application/json", new(bytes.Buffer))
	if err != nil {
		c.Print.Error(err)
		return res, err
	}

	return res, nil
}

func AppRestart(c *AppCommand, appInfo *AppInfo) (*http.Response, error) {

	c.Print.Infof("Restarting %s ...", appInfo.Name)

	res, err := utils.ExecuteHTTPRequest(fmt.Sprintf("%s/app/%s/restart", c.Host.URL, appInfo.Name), "POST", "application/json", new(bytes.Buffer))
	if err != nil {
		c.Print.Error(err)
		return res, err
	}

	return res, nil
}

func AppStop(c *AppCommand, appInfo *AppInfo) (*http.Response, error) {

	c.Print.Infof("Stopping %s ...", appInfo.Name)

	res, err := utils.ExecuteHTTPRequest(fmt.Sprintf("%s/app/%s/stop", c.Host.URL, appInfo.Name), "POST", "application/json", new(bytes.Buffer))
	if err != nil {
		c.Print.Error(err)
		return res, err
	}

	return res, nil
}

func AppRemove(c *AppCommand, appInfo *AppInfo) (*http.Response, error) {

	c.Print.Infof("Removing %s ...", appInfo.Name)

	res, err := utils.ExecuteHTTPRequest(fmt.Sprintf("%s/app/%s", c.Host.URL, appInfo.Name), "DELETE", "application/json", new(bytes.Buffer))
	if err != nil {
		c.Print.Error(err)
		return res, err
	}

	os.Remove(fmt.Sprintf("%s/%s_map", c.Paths.Storage, c.Host.Name))
	os.Remove(fmt.Sprintf("%s/%s.yaml", c.Paths.Storage, c.Host.Name))

	return res, nil
}

func (c *AppCommand) Help() string {
	return ""
}

func (c *AppCommand) Synopsis() string {
	return ""
}
