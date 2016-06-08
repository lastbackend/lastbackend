package handlers

import (
	"fmt"
	"github.com/deployithq/deployit/drivers/bolt"
	"github.com/deployithq/deployit/drivers/interfaces"
	"github.com/deployithq/deployit/drivers/log"
	"github.com/deployithq/deployit/env"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type AppInfo struct {
	UUID string `yaml:"uuid"`
	Name string `yaml:"name"`
	Tag  string `yaml:"tag"`
	URL  string `yaml:"url"`
}

var Debug bool
var SSL bool
var Host string
var Tag string
var Log bool
var Force bool
var Port int

func NewEnv() *env.Env {

	var err error

	env := &env.Env{
		Log: &log.Log{
			Logger: log.New(),
		},
		LogMode: Log,
		Port:    Port,
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

	database := bolt.Open(env.Log, fmt.Sprintf("%s/%s_map", env.StoragePath, env.Host))

	env.Storage = &bolt.Bolt{
		DB: database,
	}

	return env
}

func (a *AppInfo) Read(log interfaces.ILog, path, host string) error {
	log.Debug("Reading app info from file: ", path)

	appInfoFile, err := os.Open(fmt.Sprintf("%s/.dit/%s.yaml", path, host))
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(fmt.Sprintf("%s/.dit/%s.yaml", path, host))
			if err != nil {
				log.Error(err)
				return err
			}
			defer file.Close()
			return nil
		}
		log.Error(err)
		return err
	}
	defer appInfoFile.Close()

	data, err := ioutil.ReadAll(appInfoFile)
	if err != nil {
		log.Error(err)
		return err
	}

	err = yaml.Unmarshal(data, a)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil

}

func (a AppInfo) Write(log interfaces.ILog, path, host, uuid, name, tag, url string) error {
	log.Debug("Writing app info to file: ", path)

	appInfo, err := yaml.Marshal(AppInfo{
		UUID: uuid,
		Tag:  Tag,
		Name: name,
		URL:  url,
	})

	if err != nil {
		log.Error(err)
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/.dit/%s.yaml", path, host), appInfo, 0644)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil

}
