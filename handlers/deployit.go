package handlers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/deployithq/deployit/drivers/interfaces"
	"github.com/deployithq/deployit/env"
	"github.com/deployithq/deployit/utils"
	"gopkg.in/urfave/cli.v2"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Main deploy it handler
// - Archive all files which is in folder
// - Send it to server

func DeployIt(c *cli.Context) error {

	env := NewEnv()

	env.Log.Debug("Deploy it")

	var archiveName string = "tar.gz"
	var archivePath string = fmt.Sprintf("%s/.dit/%s", env.Path, archiveName)

	// Creating archive
	fw, err := os.Create(archivePath)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	gw := gzip.NewWriter(fw)
	tw := tar.NewWriter(gw)

	deletedFiles := []string{}

	// Listing all files from database to know what files were deleted from previous run
	// TODO Optimize logic: map stored files and check it in future
	storedFiles, err := env.Storage.ListAllFiles(env.Log)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	for _, k := range storedFiles {
		_, err := os.Stat(k)
		if err != nil {
			if os.IsNotExist(err) {
				// Gathering deleted files into slice
				deletedFiles = append(deletedFiles, k)
				err = env.Storage.Delete(env.Log, k)
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

	if err := PackFiles(env, tw, env.Path); err != nil {
		return err
	}

	tw.Close()
	gw.Close()
	fw.Close()

	bodyBuffer := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuffer)

	// Adding deleted files to request
	if len(deletedFiles) > 0 {
		delFiles, err := json.Marshal(deletedFiles)
		if err != nil {
			env.Log.Error(err)
			return err
		}

		bodyWriter.WriteField("deleted", string(delFiles))
	}

	// Adding application info to request
	bodyWriter.WriteField("name", AppName)
	bodyWriter.WriteField("tag", Tag)

	archiveInfo, err := os.Stat(archivePath)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	// If archive size is 32 it means that it is empty and we don't need to send it
	if archiveInfo.Size() != 32 {
		fh, err := os.Open(archivePath)
		if err != nil {
			env.Log.Error(err)
			return err
		}

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

		fh.Close()
	}

	bodyWriter.Close()

	env.Log.Errorf("%s/app/deploy", Host)

	// Creating response for file uploading with fields
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/app/deploy", Host), bodyBuffer)
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

	// Reading response from server
	resp_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		env.Log.Error(err)
		return err
	}
	env.Log.Debug(res.Status)
	env.Log.Debug(string(resp_body))

	res.Body.Close()

	// Deleting files
	err = os.Remove(archivePath)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	return nil
}

func PackFiles(env *env.Env, tw *tar.Writer, filesPath string) error {

	// Opening directory with files
	dir, err := os.Open(filesPath)

	if err != nil {
		env.Log.Error(err)
		return err
	}

	// Reading all files
	files, err := dir.Readdir(0)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	for _, file := range files {

		fileName := file.Name()

		// TODO Create exclude lib
		// TODO Parse .gitignore and exclude files from

		currentFilePath := fmt.Sprintf("%s/%s", filesPath, fileName)

		// Ignoring files which is not needed for build to make archive smaller
		// TODO: create base .ignore file on first application creation
		if fileName == ".git" || fileName == ".idea" || fileName == ".dit" || fileName == "node_modules" {
			continue
		}

		// If it was directory - calling this function again
		// In other case adding file to archive
		// TODO: refactor this code: after isdir checing you can call continue and do not need this large if.
		if file.IsDir() {

			if err := PackFiles(env, tw, currentFilePath); err != nil {
				return err
			}

		} else {

			// Creating path, which will be inside of archive
			newPath := strings.Replace(currentFilePath, env.Path, "", 1)[1:]

			// Creating hash
			hash := utils.Hash(fmt.Sprintf("%s:%s:%s", file.Name(), strconv.FormatInt(file.Size(), 10), file.ModTime()))

			// Reading previous hash of this file
			value, err := env.Storage.Read(env.Log, currentFilePath)
			if err != nil && err != interfaces.ErrBucketNotFound {
				return err
			}

			// If hashes are equal - add file to archive
			// TODO: if hash is equeal to value just continue. Do not necessary put any if clause
			if string(value) != hash {
				env.Log.Debug("Packing file: ", currentFilePath)

				err = env.Storage.Write(env.Log, currentFilePath, hash)
				if err != nil {
					return err
				}

				fr, err := os.Open(currentFilePath)
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

	}

	dir.Close()

	return nil

}
