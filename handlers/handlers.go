package handlers

import (
	"fmt"
	"github.com/deployithq/deployit/drivers/log"
	"github.com/deployithq/deployit/drivers/storage"
	"github.com/deployithq/deployit/env"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type AppInfo struct {
	Name string `yaml:"name"`
	Tag  string `yaml:"tag"`
	URL  string `yaml:"url"`
}

var Debug bool
var SSL bool
var Host string
var Port int

var CoreServices []string = []string{"redis"}

func NewEnv() *env.Env {

	var err error

	env := &env.Env{
		Log: &log.Log{
			Logger: log.New(),
		},
		Port: Port,
	}

	if Debug {
		env.Log.SetDebugLevel()
		env.Log.Debug("Debug mode enabled")
	}

	env.Host = Host
	env.HostUrl = fmt.Sprintf("http://%s%s", env.Host, ":"+strconv.Itoa(Port))

	if SSL {
		env.HostUrl = fmt.Sprintf("https://%s", env.Host)
	}

	env.Path, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		env.Log.Fatal(err)
	}

	env.StoragePath = fmt.Sprintf("%s/.dit", env.Path)

	err = os.Mkdir(env.StoragePath, 0766)
	if err != nil && os.IsNotExist(err) {
		env.Log.Fatal(err)
	}

	env.Storage, err = storage.Open(fmt.Sprintf("%s/%s_map", env.StoragePath, env.Host))

	return env
}

func (a *AppInfo) Read(path, host string) error {

	appInfoFile, err := os.Open(fmt.Sprintf("%s/.dit/%s.yaml", path, host))
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(fmt.Sprintf("%s/.dit/%s.yaml", path, host))
			if err != nil {
				return err
			}
			defer file.Close()
			return nil
		}
		return err
	}
	defer appInfoFile.Close()

	data, err := ioutil.ReadAll(appInfoFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, a)
	if err != nil {
		return err
	}

	return nil

}

func (a AppInfo) Write(path, host, name, tag, url string) error {

	appInfo, err := yaml.Marshal(AppInfo{
		Tag:  tag,
		Name: name,
		URL:  url,
	})

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/.dit/%s.yaml", path, host), appInfo, 0644)
	if err != nil {
		return err
	}

	return nil

}
