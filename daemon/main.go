package daemon

import (
	"github.com/codegangsta/cli"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/daemon/routes"
	"github.com/deployithq/deployit/drivers/log"
	"github.com/gorilla/mux"
	"net/http"
)

var Host string
var Debug bool

func Init(c *cli.Context) error {

	env := &env.Env{
		Log: &log.Log{
			Logger: log.New(),
		},
		Host: Host,
	}

	if Debug {
		env.Log.SetDebugLevel()
	}

	r := mux.NewRouter()

	r.HandleFunc("/", Handle(Handler{env, routes.DeployAppHandler})).Methods("POST")

	if err := http.ListenAndServe(":3000", r); err != nil {
		env.Log.Fatal(err)
	}

	return nil
}
