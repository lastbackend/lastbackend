package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/daemon/modules/app"
	"github.com/deployithq/deployit/drivers/interfaces"
	"github.com/deployithq/deployit/utils"
	"io"
	"net/http"
	"os"
)

func DeployAppHandler(env *env.Env, w http.ResponseWriter, r *http.Request) error {
	env.Log.Debug("Start uploading")

	mr, err := r.MultipartReader()
	if err != nil {
		env.Log.Error(err)
		return err
	}

	paths := []string{
		fmt.Sprintf("%s/apps", app.Default_root_path),
		fmt.Sprintf("%s/tmp", app.Default_root_path),
	}

	utils.CreateDirs(paths)

	length := r.ContentLength

	var name string = r.Header.Get("name") //r.Header.Get("x-deployit-name")
	var tag string
	var excludes []string

	a := app.App{}
	if err := a.Get(env, ""); err != nil {
		env.Log.Error(err)
		return err
	}

	if a.UUID != "" {
		a.Create(env, name, tag)
	}

	w.Header().Set("x-deployit-id", a.UUID)
	w.Header().Set("x-deployit-url", "=)")

	id := utils.GenerateID()
	var tmp_path string = fmt.Sprintf("%s/tmp/%s-tmp", app.Default_root_path, id)

	for {

		part, err := mr.NextPart()

		if err == io.EOF || part == nil {
			env.Log.Debug("Done!")
			break
		}

		if part.FormName() == "deleted" {
			env.Log.Debug("DELETE")

			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			env.Log.Debug("delete is: ", buf.String())

			if err := json.Unmarshal(buf.Bytes(), &excludes); err != nil {
				env.Log.Error(err)
				return err
			}

			continue
		}

		if part.FormName() == "tag" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			env.Log.Debug("tag is: ", buf.String())
			tag = buf.String()
			continue
		}

		if part.FormName() == "file" {

			var read int64
			var p float32

			dst, err := os.OpenFile(tmp_path, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				env.Log.Error(err)
				return err
			}

			env.Log.Debugf("Uploading progress %v%%", 0)
			w.Write([]byte(fmt.Sprintf("\rUploading progress %v%%\r", 0)))

			for {
				buffer := make([]byte, 100000)
				cBytes, err := part.Read(buffer)

				if err == io.EOF {
					env.Log.Debug("Last buffer read")
					break
				}

				read = read + int64(cBytes)

				if read <= 0 {
					break
				}

				p = float32(read*100) / float32(length)
				dst.Write(buffer[0:cBytes])

				env.Log.Debugf("Uploading progress %v", p)
				w.Write([]byte(fmt.Sprintf("\rUploading progress %v%%\r", p)))
			}

			continue
		}
	}

	env.Log.Debugf("Uploading progress %v", 100)
	w.Write([]byte(fmt.Sprintf("\rUploading progress %v%%\r", 100)))

	env.Log.Debug("incomming data info", name, tag, excludes)

	update_path := fmt.Sprintf("%s/tmp/%s", app.Default_root_path, id)
	if err := utils.Ungzip(tmp_path, update_path); err != nil {
		env.Log.Error(err)
		return err
	}

	if a.Layer != "" {
		source_path := fmt.Sprintf("%s/apps/%s", app.Default_root_path, a.Layer)
		target_path := fmt.Sprintf("%s/apps/%s", app.Default_root_path, a.Layer)

		if err := utils.ModifyLayer(source_path, update_path, target_path, excludes); err != nil {
			env.Log.Error(err)
			return err
		}
	} else {
		layer := utils.GenerateID()
		target_path := fmt.Sprintf("%s/apps/%s", app.Default_root_path, layer)

		if err := utils.CreateLayer(update_path, target_path); err != nil {
			env.Log.Error(err)
			return err
		}
	}

	if err := utils.RemoveDirs([]string{tmp_path, update_path}); err != nil {
		env.Log.Error(err)
		return err
	}

	if err := a.Deploy(env, w); err != nil {
		env.Log.Error(err)
		return err
	}

	if err := a.Start(env); err != nil {
		env.Log.Error(err)
		return err
	}

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
