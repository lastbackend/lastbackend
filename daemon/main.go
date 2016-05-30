package daemon

import (
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/daemon/routes"
	"github.com/deployithq/deployit/drivers/log"
	"github.com/gorilla/mux"
	"gopkg.in/urfave/cli.v2"
	"net/http"
	"strconv"
)

var Host string
var Port int
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
		env.Log.Debug("Debug mode enabled")
	}

	r := mux.NewRouter()

	r.HandleFunc("/app/deploy", Handle(Handler{env, routes.DeployAppHandler})).Methods("POST")

	if err := http.ListenAndServe(":"+strconv.Itoa(Port), r); err != nil {
		env.Log.Fatal("ListenAndServe: ", err)
	}

	env.Log.Debugf("Listenning... on %v port", strconv.Itoa(Port))

	return nil
}
