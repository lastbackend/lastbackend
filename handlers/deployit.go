package handlers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/deployithq/deployit/drivers/db"
	"github.com/deployithq/deployit/drivers/interfaces"
	"github.com/deployithq/deployit/drivers/log"
	"github.com/deployithq/deployit/utils"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func DeployIt(c *cli.Context) error {

	env := new(interfaces.Env)

	env.Log = &log.Log{
		Logger: log.New(),
	}

	if Debug {
		env.Log.SetDebugLevel()
	}

	currentPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	storagePath := fmt.Sprintf("%s/.dit", currentPath)

	err := os.Mkdir(storagePath, 0766)
	if err != nil && os.IsNotExist(err) {
		env.Log.Error(err)
	}

	database := db.Open(env.Log, fmt.Sprintf("%s/map", storagePath))
	defer database.Close()

	env.Database = &db.Bolt{
		DB: database,
	}

	env.Log.Debug("Deploy it")

	var archiveName string = "tar.gz"
	var pathToArchive string = fmt.Sprintf("%s/.dit/%s", currentPath, archiveName)

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
				err = env.Database.Delete(env.Log, []byte(k))
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

	err = PackFiles(env, tw, currentPath)
	if err != nil {
		return err
	}

	res, err := UploadFile(env.Log, pathToArchive, archiveName, c.String("name"), c.String("tag"))
	if err != nil {
		return err
	}
	defer res.Body.Close()

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
	defer dir.Close()

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
			err = PackFiles(env, tw, currentPath)
		} else {
			err = WriteToTarGZ(env, currentPath, tw, file)
		}

		if err != nil {
			return err
		}

	}

	return nil

}

func WriteToTarGZ(env *interfaces.Env, pathToFile string, tw *tar.Writer, file os.FileInfo) error {
	env.Log.Debug("Adding to tar.gz and hash table: " + pathToFile)

	absPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	newPath := strings.Replace(pathToFile, absPath, "", 1)[1:]

	hashData := fmt.Sprintf("%s:%s:%s", file.Name(), strconv.FormatInt(file.Size(), 10), file.ModTime())
	hash := utils.Hash([]byte(hashData))

	value, err := env.Database.Read(env.Log, []byte(pathToFile))
	if err != nil && err.Error() != "BUCKET_NOT_FOUND" {
		return err
	}

	if string(value) == hash {
		return nil
	}

	err = env.Database.Write(env.Log, []byte(pathToFile), []byte(hash))
	if err != nil {
		return err
	}

	fr, err := os.Open(pathToFile)
	if err != nil {
		env.Log.Error(err)
		return err
	}
	defer fr.Close()

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

	return nil
}

func UploadFile(log interfaces.Log, filePath, fileName, name, tag string) (*http.Response, error) {
	log.Debug("Uploading file: " + filePath)

	res := new(http.Response)

	fh, err := os.Open(filePath)
	if err != nil {
		log.Error(err)
		return res, err
	}
	defer fh.Close()

	bodyBuf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("tar.gz", filePath)
	if err != nil {
		log.Error(err)
		return res, err
	}
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		log.Error(err)
		return res, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	Body := struct {
		Name string `json:"name"`
		Tag  string `json:"tag"`
	}{name, tag}

	err = json.NewEncoder(bodyBuf).Encode(Body)
	if err != nil {
		log.Error(err)
		return &http.Response{}, err
	}

	res, err = http.Post(fmt.Sprintf("%s/app/deploy", Host), contentType, bodyBuf)
	if err != nil {
		log.Error(err)
		return res, err
	}

	return res, err
}
