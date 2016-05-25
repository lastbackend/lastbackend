package handlers

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/deployithq/deployit/drivers/db"
	"github.com/deployithq/deployit/drivers/interfaces"
	"github.com/deployithq/deployit/drivers/log"
	"os"
	"path/filepath"
)

func DeployURL(c *cli.Context) error {
	env := new(interfaces.Env)

	env.Log = &log.Log{
		Logger: log.New(),
	}

	if Debug {
		env.Log.SetDebugLevel()
	}

	currentPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	storagePath := fmt.Sprintf("%s/.dit", currentPath)

	err := os.Mkdir(storagePath, 0766)
	if err != nil && os.IsNotExist(err) {
		env.Log.Error(err)
	}

	database := db.Open(env.Log, fmt.Sprintf("%s/map", storagePath))
	defer database.Close()

	env.Database = &db.Bolt{
		DB: database,
	}

	env.Log.Debug("Deploy url")

	var url string

	if url == "" {
		err := errors.New("Empty url")
		env.Log.Error(err)
		return err
	}

	return nil
}
