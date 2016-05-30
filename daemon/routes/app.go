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

	var name, tag string
	var excludes []string

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

		if part.FormName() == "name" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			env.Log.Debug("name is: ", buf.String())
			name = buf.String()
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

	a := app.App{}
	if err := a.Get(""); err != nil {
		env.Log.Error(err)
		return err
	}

	if a.UUID != "" {
		source_path := fmt.Sprintf("%s/app/%s", app.Default_root_path, a.Layer)
		target_path := fmt.Sprintf("%s/app/%s", app.Default_root_path, a.Layer)

		if err := utils.ModifyLayer(source_path, update_path, target_path, excludes); err != nil {
			env.Log.Error(err)
			return err
		}

	} else {
		target_path := fmt.Sprintf("%s/app/%s", app.Default_root_path, utils.GenerateID())

		if err := utils.CreateLayer(update_path, target_path); err != nil {
			env.Log.Error(err)
			return err
		}
	}

	env.Log.Debug("tmp_path", tmp_path)
	env.Log.Debug("update_path", update_path)
	if err := utils.RemoveDirs([]string{tmp_path, update_path}); err != nil {
		env.Log.Error(err)
		return err
	}

	//reader, err := os.Open("temp.tar")
	//if err != nil {
	//	env.Log.Error(err)
	//	return err
	//}
	//defer reader.Close()
	//
	//or, ow := io.Pipe()
	//opts := interfaces.BuildImageOptions{
	//	Name:           "pacman:" + tag,
	//	RmTmpContainer: true,
	//	InputStream:    reader,
	//	OutputStream:   ow,
	//	RawJSONStream:  true,
	//}
	//
	//ch := make(chan error, 1)
	//
	//env.Log.Debug(">> Build <<")

	//go func() {
	//	defer ow.Close()
	//	defer close(ch)
	//	if err := env.Containers.BuildImage(opts); err != nil {
	//		env.Log.Error(err)
	//		return
	//	}
	//}()
	//
	//jsonmessage.DisplayJSONMessagesStream(or, w, os.Stdout.Fd(), term.IsTerminal(os.Stdout.Fd()), nil)
	//if err, ok := <-ch; ok {
	//	if err != nil {
	//		env.Log.Error(err)
	//		return err
	//	}
	//}

	//log.Debug(">> StartContainer <<")
	//if err := route.Context.Adapter.StartContainer(&interfaces.Container{
	//	CID: ``,
	//	Config: interfaces.Config{
	//		Image: "pacman:" + tag,
	//	},
	//	HostConfig: interfaces.HostConfig{},
	//}); err != nil {
	//	log.Error(err)
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
