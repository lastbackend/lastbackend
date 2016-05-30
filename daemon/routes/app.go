package routes

import (
	"bytes"
	"fmt"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/drivers/interfaces"
	"github.com/deployithq/deployit/utils"
		"github.com/satori/go.uuid"
	"io"
	"net/http"
	"os"
	"strconv"
	"path/filepath"
	"strings"
	"encoding/json"
	"github.com/lastbackend/scheduler/drivers/log"
	"archive/tar"
)

func Test(env *env.Env, w http.ResponseWriter, r *http.Request) error {
	for i := 0; i < 20000; i++ {
		write(env.Log, w, []byte("ping "+strconv.Itoa(i)+ "\r"))
	}

	return nil
}

func DeployAppHandler(env *env.Env, w http.ResponseWriter, r *http.Request) error {
	env.Log.Debug("Start uploading")

	mr, err := r.MultipartReader()
	if err != nil {
		env.Log.Error(err)
		return err
	}

	length := r.ContentLength

	u := uuid.NewV4()
	env.Log.Debug(u.String())

	var filename, tag string
	var excludes []string

	for {

		part, err := mr.NextPart()

		if err == io.EOF || part == nil {
			env.Log.Debug("Done!")
			break
		}

		env.Log.Debug(">> ", part.FormName())

		if part.FormName() == "delete" {
			env.Log.Debug("DELETE")

			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			env.Log.Debug("delete is: ", buf.String())

			if err:=json.Unmarshal(buf.Bytes(), &excludes); err != nil {
				env.Log.Error(err)
				return err
			}

			continue
		}

		if part.FormName() == "name" {
			env.Log.Debug("NAME")

			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			env.Log.Debug("name is: ", buf.String())
			filename = buf.String()
			continue
		}

		if part.FormName() == "tag" {
			env.Log.Debug("TAG")

			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			env.Log.Debug("tag is: ", buf.String())
			tag = buf.String()
			continue
		}

		if part.FormName() == "file" {
			env.Log.Debug("LOAD FILE")

			var read int64
			var p float32

			dst, err := os.OpenFile("upload.tar.gz", os.O_WRONLY|os.O_CREATE, 0666)

			if err != nil {
				env.Log.Error(err)
				return err
			}

			env.Log.Debugf("Uploading progress %v%%", 0)
			write(env.Log, w, []byte(fmt.Sprintf("Uploading progress %v%%", 0)))

			for {
				env.Log.Debug("Read bufer")

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

				write(env.Log, w, []byte(fmt.Sprintf("\rUploading progress %v%%", p)))
				dst.Write(buffer[0:cBytes])
			}

			continue
		}
	}

	env.Log.Debug(filename, tag)

	env.Log.Debugf("Uploading progress %v", 100)
	write(env.Log, w, []byte(fmt.Sprintf("\rUploading progress %v%%\r\r", 100)))

	if err := utils.Ungzip(env.Log, "/tmp/upload.tar.gz", "temp.tar"); err != nil {
		env.Log.Error(err)
		return err
	}

	CreateLayer(env.Log, "temp.tar", excludes)

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

func write(log interfaces.ILog, w http.ResponseWriter, data []byte)  {
	log.Debug("Flusher");

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	} else {
		log.Debug("Damn, no flush");
	}

	w.Write(data)
	w.Write([]byte("\r"))
}

func CreateLayer(log interfaces.ILog, path string, excludes []string) error {
	log.Debug("Create Layer")

	log.Debug(path)
	sources, err := os.Open(path)
	if err != nil {
		log.Error("error Open", err)
		return err
	}

	tarR2 := tar.NewReader(sources)

	target_filename := filepath.Base("3")
	target := filepath.Join("/home/unloop/Downloads/", fmt.Sprintf("%s.tar", target_filename))

	log.Debug(target)
	n, err := os.Create(target)
	if err != nil {
		log.Error("error Create", err)
		return err
	}

	tarW2 := tar.NewWriter(n)

	i := 0

	for {
		header, err := tarR2.Next()

		if err == io.EOF {

			tarW2.Flush()
			tarW2.Close()
			n.Close()

			sources.Close()

			break
		} else if err != nil {
			log.Error("error header", err)
			return err
		} else if header.Size > 1e6 {
			log.Error("huge claimed file size")
		}

		path := header.Name

		log.Debug("FLAG: ", header.Typeflag)

		switch header.Typeflag {
		case tar.TypeDir:
			// strings.TrimRight(path, "/")
			log.Debug("(", i, ")", "Dir: ", trimSuffix(path, "/"), header.Size)

			// strings.TrimRight(path, "/")
			if !include(excludes, trimSuffix(path, "/")) {
				copy(header, tarW2, tarR2)
			}

			continue
		case tar.TypeReg:
			log.Debug("(", i, ")", "Name: ", path, header.Size)

			if !include(excludes, path) {
				copy(header, tarW2, tarR2)
			}

		default:
			log.Debug("Can't save ", header.Typeflag, "in file", path)
		}

		i++
	}

	return nil
}

func copy(hdr *tar.Header, wr *tar.Writer, src *tar.Reader) error {
	// write the header to the tarball archive
	if err := wr.WriteHeader(hdr); err != nil {
		log.Error("error WriteHeader", err)
		return err
	}

	// copy the file data to the tarball
	if _, err := io.Copy(wr, src); err != nil {
		log.Error("error Copy", err)
		return err
	}

	return nil
}

func merge(hdr *tar.Header, wr *tar.Writer, src *tar.Reader) error {

	return nil
}

// Returns `true` if the target string t is in the
// slice.
func include(vs []string, t string) bool {
	return index(vs, t) >= 0
}

// Returns the first index of the target string `t`, or
// -1 if no match is found.
func index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}

	return -1
}

func trimSuffix(s, suffix string) string {

	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}