package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/deployithq/deployit/daemon/app"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/drivers/interfaces"
	"github.com/deployithq/deployit/utils"
	"io"
	"net/http"
	"os"
)

func DeployAppHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	e.Log.Debug("Start uploading")

	mr, err := r.MultipartReader()
	if err != nil {
		e.Log.Error(err)
		return err
	}

	length := r.ContentLength

	var name, tag string
	var root_path string = env.Default_root_path
	var excludes []string

	id := utils.GenerateID()
	var tmp_path string = fmt.Sprintf("%s/tmp/%s-tmp", root_path, id)

	for {

		part, err := mr.NextPart()

		if err == io.EOF || part == nil {
			e.Log.Debug("Done!")
			break
		}

		if part.FormName() == "deleted" {
			e.Log.Debug("DELETE")

			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			e.Log.Debug("delete is: ", buf.String())

			if err := json.Unmarshal(buf.Bytes(), &excludes); err != nil {
				e.Log.Error(err)
				return err
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

			dst, err := os.OpenFile(tmp_path, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				e.Log.Error(err)
				return err
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

	update_path := fmt.Sprintf("%s/tmp/%s", root_path, id)
	if err := utils.Ungzip(tmp_path, update_path); err != nil {
		e.Log.Error(err)
		return err
	}

	a := app.App{}

	if err := a.Get(e, "11dc2f82-3728-4a41-bdc3-151f90b04aac"); err != nil {
		e.Log.Error(err)
		return err
	}

	e.Log.Info("Found: ", a.UUID, a.Name)

	if a.UUID == "" {
		a.Create(e, name, tag)
	}

	w.Header().Set("x-deployit-id", a.UUID)
	w.Header().Set("x-deployit-url", "=)")

	layer := utils.GenerateID()

	if a.Layer != "" {
		source_path := fmt.Sprintf("%s/apps/%s", root_path, a.Layer)
		target_path := fmt.Sprintf("%s/apps/%s", root_path, layer)

		if err := utils.ModifyLayer(source_path, update_path, target_path, excludes); err != nil {
			e.Log.Error(err)
			return err
		}
	} else {
		target_path := fmt.Sprintf("%s/apps/%s", root_path, layer)

		if err := utils.CreateLayer(update_path, target_path); err != nil {
			e.Log.Error(err)
			return err
		}
	}

	if err := utils.RemoveDirs([]string{tmp_path, update_path}); err != nil {
		e.Log.Error(err)
		return err
	}

	//if err := a.Deploy(e); err != nil {
	//	e.Log.Error(err)
	//	return err
	//}
	//
	//if err := a.Start(e); err != nil {
	//	e.Log.Error(err)
	//	return err
	//}

	return nil
}

func write(log interfaces.ILog, w http.ResponseWriter, data []byte) {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	} else {
		log.Debug("Damn, no flush")
	}

	w.Write(data)
	w.Write([]byte("\r"))
}
