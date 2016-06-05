package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/deployithq/deployit/daemon/app"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/utils"
	"io"
	"net/http"
	"os"
	"github.com/deployithq/deployit/errors"
)

func DeployAppHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {

	e.Log.Debug("Start uploading")

	mr, err := r.MultipartReader()
	if err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	length := r.ContentLength
	id := utils.GenerateID()

	var excludes []string
	var url, uuid, name, tag string
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

		if part.FormName() == "name" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			e.Log.Debug("name is: ", buf.String())
			name = buf.String()
			continue
		}

		if part.FormName() == "tag" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			e.Log.Debug("tag is: ", buf.String())
			tag = buf.String()
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
			w.Write([]byte(fmt.Sprintf("\rUploading progress %v%%\r", 0)))

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
				w.Write([]byte(fmt.Sprintf("\rUploading progress %v%%\r", p)))
			}

			continue
		}
	}

	e.Log.Debugf("Uploading progress %v", 100)
	w.Write([]byte(fmt.Sprintf("\rUploading progress %v%%\r", 100)))

	e.Log.Debug("incomming data info", name, tag, excludes)

	a := app.App{}

	if uuid != "" {
		e.Log.Info("Get app", a.UUID)
		if err := a.Get(e, uuid); err != nil {
			e.Log.Error(err)
			return errors.InternalServerError()
		}
	} else {
		e.Log.Info("Create app")
		a.Create(e, name, tag)
	}

	if url == "" {
		e.Log.Info(targz_path, excludes)
		if err := a.Layer.CreateFromTarGz(targz_path, excludes); err != nil {
			e.Log.Error(err)
			return errors.InternalServerError()
		}
	} else {
		// TODO: Need clone source, create tar
		if err := a.Layer.CreateFromUrl(url); err != nil {
			e.Log.Error(err)
			return errors.InternalServerError()
		}
	}

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

	if err := a.Start(e); err != nil {
		e.Log.Error(err)
		return errors.InternalServerError()
	}

	w.Header().Set("x-deployit-id", a.UUID)
	w.Header().Set("x-deployit-url", "=)")

	w.Write([]byte(""))

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
