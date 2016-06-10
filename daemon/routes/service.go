package routes

import (
	"github.com/deployithq/deployit/daemon/service"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/errors"
	"github.com/deployithq/deployit/utils"
	"net/http"
)

func CreateServiceHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	e.Log.Info("Create service")

	name := utils.GetStringParamFromURL(`name`, r)
	e.Log.Debug("Deploy service", name)

	s := service.Service{}
	s.Create(e, name)

	w.Write([]byte(``))

	return nil
}

func DeployServiceHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {

	name := utils.GetStringParamFromURL(`name`, r)
	e.Log.Debug("Deploy service", name)

	s := service.Service{}
	if err := s.Get(e, name); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if err := s.Pull(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if len(s.Containers) > 0 {
		if err := s.Remove(e); err != nil {
			e.Log.Error(err)
			return errors.InternalServerError()
		}
	}

	if err := s.Start(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	w.Write([]byte(""))

	return nil
}

func StartServiceHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	name := utils.GetStringParamFromURL(`name`, r)
	e.Log.Debug("Start service", name)

	s := service.Service{}
	if err := s.Get(e, name); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if err := s.Start(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	return nil
}

func StopServiceHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	name := utils.GetStringParamFromURL(`name`, r)
	e.Log.Debug("Stop service", name)

	s := service.Service{}
	if err := s.Get(e, name); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if err := s.Stop(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	return nil
}

func RestartServiceHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	name := utils.GetStringParamFromURL(`name`, r)
	e.Log.Debug("Restart service", name)

	s := service.Service{}
	if err := s.Get(e, name); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if err := s.Restart(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	return nil
}

func RemoveServiceHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	name := utils.GetStringParamFromURL(`name`, r)
	e.Log.Debug("Remove service", name)

	s := service.Service{}
	if err := s.Get(e, name); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if err := s.Remove(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	return nil
}

func LogsServiceHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	name := utils.GetStringParamFromURL(`name`, r)
	e.Log.Debug("Logs service", name)


	return nil
}
