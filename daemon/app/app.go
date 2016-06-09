package app

import (
	"errors"
	"fmt"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/drivers/interfaces"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/satori/go.uuid"
	"io"
	"os"
	"strings"
)

type App struct {
	UUID       string                `json:"uuid" yaml:"uuid"`
	Name       string                `json:"name" yaml:"name"`
	Tag        string                `json:"tag" yaml:"tag"`
	Scale      int                   `json:"scale" yaml:"scale"`
	Layer      Layer                 `json:"layer" yaml:"layer"`
	Containers map[string]*Container `json:"container" yaml:"container"`
}

type Container struct {
	ID string `json:"id" yaml:"id"`
}

func (a *App) Get(e *env.Env, uuid string) error {
	e.Log.Info(`Get app`, uuid)

	if err := e.LDB.Read(uuid, a); err != nil {
		return err
	}

	return nil
}

func (a *App) Update(e *env.Env) error {
	e.Log.Info(`Update app`)

	if a.UUID == "" {
		return errors.New("application not found")
	}

	if err := e.LDB.Write(a.UUID, a); err != nil {
		return err
	}

	return nil
}

func (a *App) Create(e *env.Env, name, tag string) error {

	u := uuid.NewV4()

	a.UUID = u.String()
	a.Name = name
	a.Tag = tag
	a.Scale = 1
	a.Layer = Layer{}
	a.Containers = make(map[string]*Container)

	if err := e.LDB.Write(a.UUID, a); err != nil {
		return err
	}

	return nil
}

func (a *App) Build(e *env.Env, writer io.Writer) error {
	e.Log.Info(`Build app`)

	if a.Layer.ID == `` {
		return errors.New("layer not found")
	}

	tar_path := fmt.Sprintf("%s/apps/%s", env.Default_root_path, a.Layer.ID)

	reader, err := os.Open(tar_path)
	if err != nil {
		return err
	}
	defer reader.Close()

	or, ow := io.Pipe()
	opts := interfaces.BuildImageOptions{
		Name:           fmt.Sprintf("%s/%s:%s", env.Default_hub, a.Name, a.Tag),
		RmTmpContainer: true,
		InputStream:    reader,
		OutputStream:   ow,
		RawJSONStream:  true,
	}

	ch := make(chan error, 1)

	go func() {
		defer ow.Close()
		defer close(ch)
		if err := e.Containers.BuildImage(opts); err != nil {
			e.Log.Error(err)
			return
		}
	}()

	jsonmessage.DisplayJSONMessagesStream(or, writer, os.Stdout.Fd(), term.IsTerminal(os.Stdout.Fd()), nil)
	if err, ok := <-ch; ok {
		if err != nil {
			e.Log.Error(err)
			return err
		}
	}

	if err := a.Update(e); err != nil {
		return err
	}

	return nil
}

func (a *App) Start(e *env.Env) error {
	e.Log.Info(`Start app`)

	if a.UUID == "" {
		return errors.New("application not found")
	}

	//TODO: implement start with configs
	//TODO: implement scale

	// Run containers if exists
	for _, container := range a.Containers {
		if err := e.Containers.StartContainer(&interfaces.Container{
			CID:         container.ID,
			HostConfig: interfaces.HostConfig{},
		}); err != nil {
			e.Log.Error(err)
			return err
		}
	}

	for len(a.Containers) < a.Scale {
		c := &interfaces.Container{
			Config: interfaces.Config{
				Image: fmt.Sprintf("%s/%s:%s", env.Default_hub, a.Name, a.Tag),
			},
			HostConfig: interfaces.HostConfig{},
		}

		if err := e.Containers.StartContainer(c); err != nil {
			e.Log.Error(err)
			return err
		}

		a.Containers[c.CID] = &Container{
			ID: c.CID,
		}
	}

	if err := a.Update(e); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop(e *env.Env) error {
	e.Log.Info(`Stop app`)

	if a.UUID == "" {
		return errors.New("application not found")
	}

	for _, container := range a.Containers {

		if container.ID == "" {
			continue
		}

		if err := e.Containers.StopContainer(&interfaces.Container{
			CID: container.ID,
		}); err != nil {
			e.Log.Error(err)
			return err
		}
	}

	if err := a.Update(e); err != nil {
		return err
	}

	return nil
}

func (a *App) Restart(e *env.Env) error {
	e.Log.Info(`Restart app`)

	//TODO: implement start with configs
	//TODO: implement scale

	if err := a.Update(e); err != nil {
		return err
	}

	// Run containers if exists
	for _, container := range a.Containers {
		if err := e.Containers.RestartContainer(&interfaces.Container{
			CID:         container.ID,
			HostConfig: interfaces.HostConfig{},
		}); err != nil {
			e.Log.Error(err)
			return err
		}
	}

	for len(a.Containers) < a.Scale {
		c := &interfaces.Container{
			Config: interfaces.Config{
				Image: fmt.Sprintf("%s/%s:%s", env.Default_hub, a.Name, a.Tag),
			},
			HostConfig: interfaces.HostConfig{},
		}

		if err := e.Containers.StartContainer(c); err != nil {
			e.Log.Error(err)
			return err
		}

		a.Containers[c.CID] = &Container{
			ID: c.CID,
		}
	}

	return nil
}

func (a *App) Remove(e *env.Env) error {
	e.Log.Info(`Remove app`)

	if a.UUID == "" {
		return errors.New("application not found")
	}

	for key, container := range a.Containers {
		if container.ID != "" {
			if err := e.Containers.RemoveContainer(&interfaces.Container{
				CID: container.ID,
			}); err != nil {
				e.Log.Error(err)
				index := strings.Index(err.Error(), "No such container")
				if index != -1 {
					e.Log.Info(`Clear record in db`)
					delete(a.Containers, key)
					continue
				}

				return err
			}
		}

		delete(a.Containers, key)
	}

	if err := a.Update(e); err != nil {
		return err
	}

	return nil
}

func (a *App) Destroy(e *env.Env) error {
	e.Log.Info(`Destroy app`)

	if a.UUID == "" {
		return errors.New("application not found")
	}

	if err := a.Remove(e); err != nil {
		return err
	}

	if err := e.LDB.Remove(a.UUID); err != nil {
		return err
	}

	return nil
}
