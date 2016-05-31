package handlers

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/deployithq/deployit/env"
	"github.com/deployithq/deployit/utils"
	"github.com/fatih/color"
	"gopkg.in/urfave/cli.v2"
	"io"
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

	var archiveName string = "tar.gz"
	var archivePath string = fmt.Sprintf("%s/.dit/%s", env.Path, archiveName)

	appInfo := new(AppInfo)
	err := appInfo.Read(env.Log, env.Path, env.Host)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	if appInfo.Name == "" {
		appInfo.Name = utils.AppName(env.Path)
		appInfo.Tag = Tag
		color.Cyan("Creating app: %s", appInfo.Name)
	} else {
		color.Cyan("Updating app: %s", appInfo.Name)
	}

	// Creating archive
	fw, err := os.Create(archivePath)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	gw := gzip.NewWriter(fw)
	tw := tar.NewWriter(gw)

	// Deleting archive after function ends
	defer func() {
		env.Log.Debug("Deleting archive: ", archivePath)

		fw.Close()
		gw.Close()
		tw.Close()

		// Deleting files
		err = os.Remove(archivePath)
		if err != nil {
			env.Log.Error(err)
			return
		}
	}()

	// Listing all files from database to know what files were deleted from previous run
	storedFiles, err := env.Storage.ListAllFiles(env.Log)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	// TODO Include deleted folders to deletedFiles like "nginx/"

	color.Cyan("Packing files")
	storedFiles, err = PackFiles(env, tw, env.Path, storedFiles)
	if err != nil {
		return err
	}

	deletedFiles := []string{}

	for key, _ := range storedFiles {
		env.Log.Debug("Deleting: ", key)
		err = env.Storage.Delete(env.Log, key)
		if err != nil {
			return err
		}
		deletedFiles = append(deletedFiles, key)
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

		bodyWriter.WriteField("x-deployit-deleted", string(delFiles))
	}

	// Adding application info to request
	if appInfo.UUID == "" {
		bodyWriter.WriteField("x-deployit-name", appInfo.Name)
	} else {
		bodyWriter.WriteField("x-deployit-id", appInfo.UUID)
	}

	bodyWriter.WriteField("x-deployit-tag", appInfo.Name)

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

	env.Log.Debugf("%s/app/deploy", Host)

	// Creating response for file uploading with fields
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/app/deploy", env.HostUrl), bodyBuffer)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	color.Cyan("Uploading sources")

	// TODO Show uploading progress

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	if Log {
		color.Cyan("Logs: ")
		reader := bufio.NewReader(res.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				break
			}
			fmt.Println(string(line))
		}
	}

	appInfo.URL = res.Header.Get("x-deployit-url")

	color.Cyan(appInfo.URL)

	// TODO Handle errors from http - clear DB if was first run

	if appInfo.UUID == "" {
		appInfo.UUID = res.Header.Get("x-deployit-id")
	}

	err = appInfo.Write(env.Log, env.Path, env.Host, appInfo.UUID, appInfo.Name, appInfo.Tag, appInfo.URL)
	if err != nil {
		env.Log.Error(err)
		return err
	}

	res.Body.Close()

	color.Cyan("Done")

	return nil
}

func PackFiles(env *env.Env, tw *tar.Writer, filesPath string, storedFiles map[string]string) (map[string]string, error) {

	// Opening directory with files
	dir, err := os.Open(filesPath)

	if err != nil {
		env.Log.Error(err)
		return storedFiles, err
	}

	// Reading all files
	files, err := dir.Readdir(0)
	if err != nil {
		env.Log.Error(err)
		return storedFiles, err
	}

	for _, file := range files {

		fileName := file.Name()

		currentFilePath := fmt.Sprintf("%s/%s", filesPath, fileName)

		// Ignoring files which is not needed for build to make archive smaller
		// TODO: create base .ignore file on first application creation
		// TODO Create exclude lib
		// TODO Parse .gitignore and exclude files from
		if fileName == ".git" || fileName == ".idea" || fileName == ".dit" || fileName == "node_modules" {
			continue
		}

		// If it was directory - calling this function again
		// In other case adding file to archive
		if file.IsDir() {
			storedFiles, err = PackFiles(env, tw, currentFilePath, storedFiles)
			if err != nil {
				return storedFiles, err
			}
			continue
		}

		// Creating path, which will be inside of archive
		newPath := strings.Replace(currentFilePath, env.Path, "", 1)[1:]

		// Creating hash
		hash := utils.Hash(fmt.Sprintf("%s:%s:%s", file.Name(), strconv.FormatInt(file.Size(), 10), file.ModTime()))

		if storedFiles[newPath] == hash {
			delete(storedFiles, newPath)
			continue
		}

		delete(storedFiles, newPath)

		// If hashes are not equal - add file to archive
		env.Log.Debug("Packing file: ", currentFilePath)

		err = env.Storage.Write(env.Log, newPath, hash)
		if err != nil {
			return storedFiles, err
		}

		fr, err := os.Open(currentFilePath)
		if err != nil {
			env.Log.Error(err)
			return storedFiles, err
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
			return storedFiles, err
		}

		_, err = io.Copy(tw, fr)
		if err != nil {
			env.Log.Error(err)
			return storedFiles, err
		}

		fr.Close()

	}

	dir.Close()

	return storedFiles, err

}
