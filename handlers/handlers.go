package handlers

import (
	"fmt"
	"github.com/deployithq/deployit/drivers/bolt"
	"github.com/deployithq/deployit/drivers/log"
	"github.com/deployithq/deployit/env"
	"os"
	"path/filepath"
)

var Debug bool
var Host string
var Tag string

func NewEnv() *env.Env {

	var err error

	env := &env.Env{
		Log: &log.Log{
			Logger: log.New(),
		},
		Host: Host,
	}

	if Debug {
		env.Log.SetDebugLevel()
		env.Log.Debug("Debug node enabled")
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

	database := bolt.Open(env.Log, fmt.Sprintf("%s/map", env.StoragePath))

	env.Storage = &bolt.Bolt{
		DB: database,
	}

	return env
}
