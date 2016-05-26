package handlers

import (
	"fmt"
	"github.com/deployithq/deployit/drivers/db"
	"github.com/deployithq/deployit/drivers/interfaces"
	"github.com/deployithq/deployit/drivers/log"
	"os"
	"path/filepath"
)

var Debug bool
var Host string
var AppName string
var Tag string

func NewEnv() *interfaces.Env {

	var err error

	env := &interfaces.Env{
		Log: &log.Log{
			Logger: log.New(),
		},
		Host: Host,
	}

	if Debug {
		env.Log.SetDebugLevel()
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

	database := db.Open(env.Log, fmt.Sprintf("%s/map", env.StoragePath))

	env.Database = &db.Bolt{
		DB: database,
	}

	return env
}
