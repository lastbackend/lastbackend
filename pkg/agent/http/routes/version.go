package routes

import (
	"net/http"
	"github.com/Sirupsen/logrus"
)

func VersionGetR(w http.ResponseWriter, r *http.Request) {

	var (
		err error
	)

	logrus.Debug("Get user handler")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("{}"))
	if err != nil {
		logrus.Error("Error: write response", err.Error())
		return
	}
}