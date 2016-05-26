package daemon

import (
	"fmt"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/daemon/routes"
	"github.com/deployithq/deployit/drivers/log"
	"github.com/gorilla/mux"
	"gopkg.in/urfave/cli.v2"
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

	fmt.Println(c.Bool("debug"))

	if Debug {
		env.Log.SetDebugLevel()
	}

	env.Log.Debug("Init")

	r := mux.NewRouter()

	r.HandleFunc("/app/deploy", Handle(Handler{env, routes.DeployAppHandler})).Methods("POST")

	if err := http.ListenAndServe(":3000", r); err != nil {
		env.Log.Fatal(err)
	}

	env.Log.Debug("Listenning... on %d port", 3000)

	return nil
}
