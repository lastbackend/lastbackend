package app

import (
	"github.com/deployithq/deployit/daemon/env"
	"github.com/satori/go.uuid"
)

type App struct {
	UUID  string `json:"uuid" yaml:"uuid"`
	Name  string `json:"name" yaml:"name"`
	Tag   string `json:"tag" yaml:"tag"`
	Layer string `json:"layer" yaml:"layer"`
}

func (a *App) Get(e *env.Env, uuid string) error {
	e.Log.Info(`Get app`, uuid)

	if err := e.LDB.Get(uuid, a); err != nil {
		return err
	}

	return nil
}

func (a *App) Create(e *env.Env, name, tag string) error {
	e.Log.Info(`Create app`, name, tag)

	u := uuid.NewV4()
	a = &App{
		UUID: u.String(),
		Name: name,
		Tag:  tag,
	}

	if err := e.LDB.Set(a.UUID, &a); err != nil {
		return err
	}

	return nil
}

func (a *App) Deploy(e *env.Env, uuid string) error {
	e.Log.Info(`Deploy app`)

	return nil
}

func (a App) Remove(e *env.Env, uuid string) error {
	e.Log.Info(`Remove app`, uuid)
	return nil
}
