package app

import ()
import (
	"errors"
	"fmt"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/drivers/interfaces"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"io"
	"net/http"
	"os"
)

type App struct {
	UUID      string
	Name      string
	Layer     string
	Tag       string
	Namespace string
}

func (a *App) Create(e *env.Env, name, tag string) error {
	return nil
}

func (a *App) Get(e *env.Env, uuid string) error {
	return nil
}

func (a *App) Deploy(e *env.Env, w http.ResponseWriter) error {

	if a.UUID == "" {
		return errors.New("application not found")
	}

	reader, err := os.Open("temp.tar")
	if err != nil {
		e.Log.Error(err)
		return err
	}
	defer reader.Close()

	or, ow := io.Pipe()
	opts := interfaces.BuildImageOptions{
		Name:           fmt.Sprintf("%s:%s", a.Name, a.Tag), //a.Namespace,
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

	jsonmessage.DisplayJSONMessagesStream(or, w, os.Stdout.Fd(), term.IsTerminal(os.Stdout.Fd()), nil)
	if err, ok := <-ch; ok {
		if err != nil {
			e.Log.Error(err)
			return err
		}
	}

	return nil
}

func (a *App) Start(e *env.Env) error {

	if a.UUID == "" {
		return errors.New("application not found")
	}

	if err := e.Containers.StartContainer(&interfaces.Container{
		CID: ``,
		Config: interfaces.Config{
			Image: fmt.Sprintf("%s:%s", a.Name, a.Tag), //a.Namespace,
		},
		HostConfig: interfaces.HostConfig{},
	}); err != nil {
		e.Log.Error(err)
		return err
	}

	return nil
}

func (a *App) Restart(e *env.Env) error {

	if a.UUID == "" {
		return errors.New("application not found")
	}

	return nil
}

func (a *App) Stop(e *env.Env) error {

	if a.UUID == "" {
		return errors.New("application not found")
	}

	return nil
}

func (a *App) Remove(e *env.Env) error {

	if a.UUID == "" {
		return errors.New("application not found")
	}

	return nil
}
