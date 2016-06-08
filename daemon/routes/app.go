package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/deployithq/deployit/daemon/app"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/errors"
	"github.com/deployithq/deployit/utils"
	"io"
	"net/http"
	"os"
)

func CreateAppHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	e.Log.Info("Create app")

	payload := struct {
		Name string `json:"name"`
		Tag  string `json:"tag"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return errors.InvalidIncomingJSON()
	}

	if payload.Name == "" {
		return errors.ParamInvalid(`name`)
	}

	if payload.Tag == "" {
		return errors.ParamInvalid(`tag`)
	}

	a := app.App{}
	a.Create(e, payload.Name, payload.Tag)

	w.Write([]byte(fmt.Sprintf(`{"uuid":"%s"}`, a.UUID)))

	return nil
}

func DeployAppHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {

	uuid := utils.GetStringParamFromURL(`id`, r)
	e.Log.Debug("Deploy app", uuid)

	mr, err := r.MultipartReader()
	if err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	length := r.ContentLength
	id := utils.GenerateID()

	var excludes []string
	var root_path string = env.Default_root_path
	var targz_path string = fmt.Sprintf("%s/tmp/%s-tmp", root_path, id)

	for {
		part, err := mr.NextPart()

		if err == io.EOF || part == nil {
			e.Log.Debug("Done!")
			break
		}

		if part.FormName() == "deleted" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			e.Log.Debug("delete is: ", buf.String())

			if err := json.Unmarshal(buf.Bytes(), &excludes); err != nil {
				e.Log.Error(err)
				return errors.InvalidIncomingJSON()
			}

			continue
		}

		if part.FormName() == "file" {

			var read int64
			var p float32

			dst, err := os.OpenFile(targz_path, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				e.Log.Error(err)
				return errors.InternalServerError()
			}

			e.Log.Debugf("Uploading progress %v%%", 0)
			// w.Write([]byte(fmt.Sprintf("\rUploading progress %v%%\r", 0)))

			for {
				buffer := make([]byte, 100000)
				cBytes, err := part.Read(buffer)

				if err == io.EOF {
					e.Log.Debug("Last buffer read")
					break
				}

				read = read + int64(cBytes)

				if read <= 0 {
					break
				}

				p = float32(read*100) / float32(length)
				dst.Write(buffer[0:cBytes])

				e.Log.Debugf("Uploading progress %v", p)
				// w.Write([]byte(fmt.Sprintf("\rUploading progress %v%%\r", p)))
			}

			continue
		}
	}

	e.Log.Debugf("Uploading progress %v", 100)
	// w.Write([]byte(fmt.Sprintf("\rUploading progress %v%%\r", 100)))

	e.Log.Debug("incomming data info", excludes)

	a := app.App{}
	if err := a.Get(e, uuid); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	e.Log.Info(targz_path, excludes)
	if err := a.Layer.CreateFromTarGz(targz_path, excludes); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	w.Header().Set("x-deployit-id", a.UUID)
	w.Header().Set("x-deployit-url", "=)")

	if err := a.Build(e, w); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if err := a.Update(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if err := utils.RemoveDirs([]string{targz_path}); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if len(a.Containers) > 0 {
		if err := a.Remove(e); err != nil {
			e.Log.Error(err)
			return errors.InternalServerError()
		}
	}

	if err := a.Start(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	w.Write([]byte(""))

	return nil
}

func StartAppHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	uuid := utils.GetStringParamFromURL(`id`, r)
	e.Log.Debug("Start app", uuid)

	a := app.App{}
	if err := a.Get(e, uuid); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if err := a.Start(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	return nil
}

func StopAppHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	uuid := utils.GetStringParamFromURL(`id`, r)
	e.Log.Debug("Stop app", uuid)

	a := app.App{}
	if err := a.Get(e, uuid); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if err := a.Stop(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	return nil
}

func RestartAppHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	uuid := utils.GetStringParamFromURL(`id`, r)
	e.Log.Debug("Restart app", uuid)

	a := app.App{}
	if err := a.Get(e, uuid); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if err := a.Restart(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	return nil
}

func RemoveAppHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	uuid := utils.GetStringParamFromURL(`id`, r)
	e.Log.Debug("Remove app", uuid)

	a := app.App{}
	if err := a.Get(e, uuid); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	if err := a.Remove(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	return nil
}

//func write(log interfaces.ILog, w http.ResponseWriter, data []byte) {
//	if f, ok := w.(http.Flusher); ok {
//		f.Flush()
//	} else {
//		log.Debug("Damn, no flush")
//	}
//
//	w.Write(data)
//	w.Write([]byte("\r"))
//}
