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
	Layer      Layer                 `json:"layer" yaml:"layer"`
	Containers map[string]*Container `json:"container" yaml:"container"`
	Config     Config                `json:"config" yaml:"config"`
}

type Container struct {
	ID string `json:"id" yaml:"id"`
}

func (a *App) Get(e *env.Env, name string) error {
	e.Log.Info(`Get app `, name)

	if err := e.LDB.Read(name, a); err != nil {
		return err
	}

	return nil
}

func (a *App) Update(e *env.Env) error {
	e.Log.Info(`Update app `, a.Name)

	if a.UUID == "" {
		return errors.New("application not found")
	}

	if err := e.LDB.Write(a.Name, a); err != nil {
		return err
	}

	return nil
}

func (a *App) Create(e *env.Env, hub, name, tag string) error {
	e.Log.Info(`Create app`)

	u := uuid.NewV4()

	a.UUID = u.String()
	a.Name = name
	a.Tag = tag
	a.Containers = make(map[string]*Container)
	a.Layer = Layer{}

	a.Config = Config{}
	a.Config.Create(e, hub, a.Name, a.Tag)

	if err := e.LDB.Write(a.Name, a); err != nil {
		return err
	}

	return nil
}

func (a *App) Build(e *env.Env, writer io.Writer) error {
	e.Log.Info(`Build app`)

	if a.Layer.ID == `` {
		return errors.New("layer not found")
	}

	if err := a.Config.Sync(e, a.Layer.ID); err != nil {
		return err
	}

	tar_path := fmt.Sprintf("%s/apps/%s", env.Default_root_path, a.Layer.ID)

	reader, err := os.Open(tar_path)
	if err != nil {
		return err
	}

	or, ow := io.Pipe()
	opts := interfaces.BuildImageOptions{
		Name:           a.Config.Image,
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

	reader.Close()
	or.Close()

	return nil
}

func (a *App) Start(e *env.Env) error {
	e.Log.Info(`Start app`)

	if a.UUID == "" {
		return errors.New("application not found")
	}

	//TODO: implement scale

	hcfg := interfaces.HostConfig{
		Memory:     a.Config.Memory,
		Ports:      a.Config.Ports,
		Binds:      a.Config.Volumes,
		Privileged: false,
		RestartPolicy: interfaces.RestartPolicyConfig{
			Attempt: 10,
			Name:    "always",
		},
	}

	// Run containers if exists
	for _, container := range a.Containers {

		if err := e.Containers.StartContainer(&interfaces.Container{
			CID:        container.ID,
			HostConfig: hcfg,
		}); err != nil {
			e.Log.Error(err)
			return err
		}
	}

	for len(a.Containers) == 0 {

		c := &interfaces.Container{
			Config: interfaces.Config{
				Image:   a.Config.Image,
				Memory:  a.Config.Memory,
				Ports:   a.Config.Ports,
				Volumes: a.Config.Volumes,
				Env:     a.Config.Env,
			},
			HostConfig: hcfg,
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

	hcfg := interfaces.HostConfig{
		Memory:     a.Config.Memory,
		Ports:      a.Config.Ports,
		Binds:      a.Config.Volumes,
		Privileged: false,
		RestartPolicy: interfaces.RestartPolicyConfig{
			Attempt: 10,
			Name:    "always",
		},
	}

	// Run containers if exists
	for _, container := range a.Containers {
		if err := e.Containers.RestartContainer(&interfaces.Container{
			CID:        container.ID,
			HostConfig: hcfg,
		}); err != nil {
			e.Log.Error(err)
			return err
		}
	}

	for len(a.Containers) == 0 {
		c := &interfaces.Container{
			Config: interfaces.Config{
				Image:   a.Config.Image,
				Memory:  a.Config.Memory,
				Ports:   a.Config.Ports,
				Volumes: a.Config.Volumes,
				Env:     a.Config.Env,
			},
			HostConfig: hcfg,
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

// Todo: temporary solution
func (a *App) Ports(e *env.Env) (int64, error) {

	var port int64

	if len(a.Containers) == 0 {
		return port, nil
	}

	for _, container := range a.Containers {
		ports, err := e.Containers.InspectContainers(&interfaces.Container{
			CID: container.ID,
		})

		if err != nil {
			e.Log.Error(err)
			return port, err
		}

		port = ports[0]

		break
	}

	return port, nil
}
