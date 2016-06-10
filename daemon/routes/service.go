package routes

import (
	"github.com/deployithq/deployit/daemon/service"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/errors"
	"github.com/deployithq/deployit/utils"
	"net/http"
	"fmt"
)

func CreateServiceHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {

	name := utils.GetStringParamFromURL(`name`, r)
	e.Log.Info("Deploy service handler ", name)

	s := service.Service{}

	s.Get(e, name)

	if s.UUID != `` {
		port, err := s.Ports(e)
		if err != nil {
			e.Log.Error(err)
			return errors.InternalServerError()
		}

		w.Write([]byte(fmt.Sprintf(`{"port":%d}`, port)))
		return nil
	}

	if err := s.Create(e, name); err != nil {
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

	port, err := s.Ports(e)
	if err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	w.Write([]byte(fmt.Sprintf(`{"port":%d}`, port)))

	return nil
}

func StartServiceHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	name := utils.GetStringParamFromURL(`name`, r)
	e.Log.Debug("Start service handler ", name)

	s := service.Service{}
	if err := s.Get(e, name); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if err := s.Start(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	port, err := s.Ports(e)
	if err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	w.Write([]byte(fmt.Sprintf(`{"port":%d}`, port)))

	return nil
}

func StopServiceHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	name := utils.GetStringParamFromURL(`name`, r)
	e.Log.Debug("Stop service handler ", name)

	s := service.Service{}
	if err := s.Get(e, name); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if err := s.Stop(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	w.Write([]byte(``))

	return nil
}

func RestartServiceHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	name := utils.GetStringParamFromURL(`name`, r)
	e.Log.Debug("Restart service handler ", name)

	s := service.Service{}
	if err := s.Get(e, name); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if err := s.Restart(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	port, err := s.Ports(e)
	if err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	w.Write([]byte(fmt.Sprintf(`{"port":%d}`, port)))

	return nil
}

func RemoveServiceHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	name := utils.GetStringParamFromURL(`name`, r)
	e.Log.Debug("Remove service handler ", name)

	s := service.Service{}
	if err := s.Get(e, name); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if err := s.Remove(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	w.Write([]byte(``))

	return nil
}

func LogsServiceHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	name := utils.GetStringParamFromURL(`name`, r)
	e.Log.Debug("Logs service handler ", name)

	w.Write([]byte(``))

	return nil
}
