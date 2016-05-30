package app

import ()

const (
	Default_root_path string = "/var/lib/deployit"
)

type App struct {
	UUID  string
	Layer string
}

func (a *App) Get(uuid string) error {
	return nil
}

func (a *App) Deploy(uuid string) error {
	return nil
}

func (a *App) Start(uuid string) error {
	return nil
}

func (a *App) Stop(uuid string) error {
	return nil
}

func (a *App) Remove(uuid string) error {
	return nil
}
