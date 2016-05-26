package handlers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/deployithq/deployit/drivers/interfaces"
	"github.com/deployithq/deployit/utils"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func DeployIt(c *cli.Context) error {

	env := NewEnv()

	env.Log.Debug("Deploy it")

	var archiveName string = "tar.gz"
	var pathToArchive string = fmt.Sprintf("%s/.dit/%s", env.CurrentPath, archiveName)

	fw, err := os.Create(pathToArchive)
	if err != nil {
		env.Log.Error(err)
		return err
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	deletedFiles := []string{}

	// TODO Optimize logic: map stored files and check it in future
	storedFiles, err := env.Database.ListAllFiles(env.Log)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	for _, k := range storedFiles {
		_, err := os.Stat(k)
		if err != nil {
			if os.IsNotExist(err) {
				deletedFiles = append(deletedFiles, k)
				err = env.Database.Delete(env.Log, k)
				if err != nil {
					env.Log.Error(err)
					return err
				}
				continue
			} else {
				env.Log.Error(err)
				return err
			}
		}
	}

	err = PackFiles(env, tw, env.CurrentPath)
	if err != nil {
		return err
	}

	tw.Close()
	gw.Close()
	fw.Close()

	bodyBuf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuf)

	if len(deletedFiles) > 0 {
		delFiles, err := json.Marshal(deletedFiles)
		if err != nil {
			env.Log.Error(err)
			return err
		}

		bodyWriter.WriteField("deleted", string(delFiles))
	}

	bodyWriter.WriteField("name", AppName)
	bodyWriter.WriteField("tag", Tag)

	archive, err := os.Stat(pathToArchive)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	if archive.Size() != 32 {
		fh, err := os.Open(pathToArchive)
		if err != nil {
			env.Log.Error(err)
			return err
		}
		defer fh.Close()

		fileWriter, err := bodyWriter.CreateFormFile("file", "tar.gz")
		if err != nil {
			env.Log.Error(err)
			return err
		}

		_, err = io.Copy(fileWriter, fh)
		if err != nil {
			env.Log.Error(err)
			return err
		}
	}

	bodyWriter.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/app/deploy", Host), bodyBuf)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	defer res.Body.Close()

	resp_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		env.Log.Error(err)
		return err
	}
	env.Log.Debug(res.Status)
	env.Log.Debug(string(resp_body))

	err = os.Remove(pathToArchive)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	return nil
}

func PackFiles(env *interfaces.Env, tw *tar.Writer, pathToFiles string) error {
	env.Log.Debug("Packing files")

	dir, err := os.Open(pathToFiles)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	files, err := dir.Readdir(0)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	for _, file := range files {

		fileName := file.Name()

		// TODO Create exclude lib
		// TODO Parse .gitignore and exclude files from

		currentPath := fmt.Sprintf("%s/%s", pathToFiles, fileName)

		if fileName == ".git" || fileName == ".idea" || fileName == ".dit" || fileName == "node_modules" {
			continue
		}

		if file.IsDir() {

			if err = PackFiles(env, tw, currentPath); err != nil {
				return err
			}

		} else {

			newPath := strings.Replace(currentPath, pathToFiles, "", 1)[1:]

			hashData := fmt.Sprintf("%s:%s:%s", file.Name(), strconv.FormatInt(file.Size(), 10), file.ModTime())
			hash := utils.Hash([]byte(hashData))

			value, err := env.Database.Read(env.Log, currentPath)
			if err != nil && err != interfaces.ErrBucketNotFound {
				return err
			}

			if string(value) == hash {
				return nil
			}

			err = env.Database.Write(env.Log, currentPath, hash)
			if err != nil {
				return err
			}

			fr, err := os.Open(currentPath)
			if err != nil {
				env.Log.Error(err)
				return err
			}

			h := &tar.Header{
				Name:    newPath,
				Size:    file.Size(),
				Mode:    int64(file.Mode()),
				ModTime: file.ModTime(),
			}

			err = tw.WriteHeader(h)
			if err != nil {
				env.Log.Error(err)
				return err
			}

			_, err = io.Copy(tw, fr)
			if err != nil {
				env.Log.Error(err)
				return err
			}

			fr.Close()

		}

	}

	dir.Close()

	return nil

}
