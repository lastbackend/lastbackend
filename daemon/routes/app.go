package routes

import (
	"bytes"
	"fmt"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/drivers/interfaces"
	"github.com/deployithq/deployit/utils"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
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

	length := r.ContentLength

	var filename, tag string

	// TODO: I guess it will be more productive to create a special header with first 10 bytes of is
	// the header size and cut this headers from incomming buffer. The main idea is to cutl tech data
	// from privided tarrball, if it's possible ofcource

	for {
		part, err := mr.NextPart()

		if err == io.EOF {
			env.Log.Debug("Done!")
			break
		}

		if part.FormName() == "delete" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			env.Log.Debug("delete is: ", buf.String())
			continue
		}

		if part.FormName() == "name" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			env.Log.Debug("name is: ", buf.String())
			filename = buf.String()
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

			dst, err := os.OpenFile("upload.tar.gz", os.O_WRONLY|os.O_CREATE, 0666)

			if err != nil {
				env.Log.Error(err)
				return err
			}

			env.Log.Debugf("Uploading progress %v%%", 0)
			w.Write([]byte(fmt.Sprintf("Uploading progress %v%%", 0)))

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

				env.Log.Debugf("Uploading progress %v", p)

				w.Write([]byte(fmt.Sprintf("\rUploading progress %v%%", p)))
				dst.Write(buffer[0:cBytes])
			}
			continue
		}
	}

	env.Log.Debug(filename, tag)

	env.Log.Debugf("Uploading progress %v", 100)
	w.Write([]byte(fmt.Sprintf("\rUploading progress %v%%\r\r", 100)))

	if err := utils.Ungzip(env.Log, "upload.tar.gz", "temp.tar"); err != nil {
		env.Log.Error(err)
		return err
	}

	reader, err := os.Open("temp.tar")
	if err != nil {
		env.Log.Error(err)
		return err
	}
	defer reader.Close()

	or, ow := io.Pipe()
	opts := interfaces.BuildImageOptions{
		Name:           "pacman:" + tag,
		RmTmpContainer: true,
		InputStream:    reader,
		OutputStream:   ow,
		RawJSONStream:  true,
	}

	ch := make(chan error, 1)

	env.Log.Debug(">> Build <<")
	go func() {
		defer ow.Close()
		defer close(ch)
		if err := env.Containers.BuildImage(opts); err != nil {
			env.Log.Error(err)
			return
		}
	}()

	jsonmessage.DisplayJSONMessagesStream(or, w, os.Stdout.Fd(), term.IsTerminal(os.Stdout.Fd()), nil)
	if err, ok := <-ch; ok {
		if err != nil {
			env.Log.Error(err)
			return err
		}
	}

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
