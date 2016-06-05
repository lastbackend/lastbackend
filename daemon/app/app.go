package app

import (
	"errors"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/satori/go.uuid"
	"fmt"
	"os"
	"io"
	"github.com/deployithq/deployit/drivers/interfaces"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
)

type App struct {
	UUID  string `json:"uuid" yaml:"uuid"`
	Name  string `json:"name" yaml:"name"`
	Tag   string `json:"tag" yaml:"tag"`
	Layer Layer  `json:"layer" yaml:"layer"`
}

func (a *App) Get(e *env.Env, uuid string) error {
	e.Log.Info(`Get app`, uuid)

	if err := e.LDB.Get(uuid, a); err != nil {
		return err
	}

	return nil
}

func (a *App) Update(e *env.Env) error {
	e.Log.Info(`Update app`)

	if a.UUID == "" {
		return errors.New("application not found")
	}

	if err := e.LDB.Set(a.UUID, a); err != nil {
		return err
	}
	return nil
}

func (a *App) Create(e *env.Env, name, tag string) error {

	u := uuid.NewV4()

	a.UUID = u.String()
	a.Name = name
	a.Tag = tag
	a.Layer = Layer{}

	if err := e.LDB.Set(a.UUID, a); err != nil {
		return err
	}

	return nil
}

func (a *App) Build(e *env.Env, writer io.Writer) error {
	e.Log.Info(`Build app`)

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

	return nil
}

func (a *App) Start(e *env.Env) error {
	e.Log.Info(`Start app`)

	if a.UUID == "" {
		return errors.New("application not found")
	}

	if err := e.Containers.StartContainer(&interfaces.Container{
		CID: ``,
		Config: interfaces.Config{
			Image: fmt.Sprintf("%s/%s:%s", env.Default_hub, a.Name, a.Tag),
		},
		HostConfig: interfaces.HostConfig{},
	}); err != nil {
		e.Log.Error(err)
		return err
	}

	return nil
}

func (a *App) Stop(e *env.Env) error {
	return nil
}

func (a *App) Restart(e *env.Env) error {
	return nil
}

func (a *App) Remove(e *env.Env) error {
	return nil
}
