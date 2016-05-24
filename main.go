package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
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
)

var Debug bool

type Env struct {
	Log      interfaces.Log
	Database interfaces.DB
}

func main() {
	app := cli.NewApp()
	app.Name = "deployit"
	app.Usage = ""

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Debug mode",
			Destination: &Debug,
		},
	}

	app.Action = Init

	app.Run(os.Args)

}

func Init(c *cli.Context) error {

	env := new(Env)

	env.Log = &log.Log{
		Logger: log.New(),
	}

	if Debug {
		env.Log.SetDebugLevel()
	}

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	storagePath := fmt.Sprintf("%s/.dit", dir)

	err := os.Mkdir(storagePath, 0766)
	if err != nil && os.IsNotExist(err) {
		env.Log.Error(err)
	}

	database := db.Open(env.Log, fmt.Sprintf("%s/map", storagePath))
	defer database.Close()

	env.Database = &db.Bolt{
		DB: database,
	}

	cmd := string(c.Args()[0])

	switch cmd {
	case "it":
		DeployIt(env, dir)
	case "url":
		DeployURL(env, "")
	}

	return nil
}

func DeployIt(env *Env, pathToFiles string) error {
	env.Log.Debug("Deploy it")

	var archiveName string = "tar.gz"

	var pathToArchive string = fmt.Sprintf("%s/.dit/%s", pathToFiles, archiveName)

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

	err = PackFiles(env, tw, pathToFiles)
	if err != nil {
		return err
	}

	res, err := UploadFile(env.Log, pathToArchive, archiveName, "")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// TODO Remove archive

	return nil

}

func DeployURL(env *Env, url string) error {
	env.Log.Debug("Deploy url")

	if url == "" {
		err := errors.New("Empty url")
		env.Log.Error(err)
		return err
	}

	return nil
}

func UploadFile(log interfaces.Log, filePath, fileName, url string) (*http.Response, error) {
	log.Debug("Uploading file: " + filePath)

	res := new(http.Response)

	bodyBuf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("code", fileName)
	if err != nil {
		log.Error(err)
		return res, err
	}

	fh, err := os.Open(filePath)
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

	res, err = http.Post("", contentType, bodyBuf)
	if err != nil {
		log.Error(err)
		return res, err
	}

	return res, err
}

func PackFiles(env *Env, tw *tar.Writer, pathToFiles string) error {
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

		if fileName == ".git" || fileName == ".idea" || fileName == ".dit" {
			continue
		}

		currentPath := fmt.Sprintf("%s/%s", pathToFiles, fileName)

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

func WriteToTarGZ(env *Env, pathToFile string, tw *tar.Writer, file os.FileInfo) error {
	env.Log.Debug("Adding to tar.gz and hash table: " + pathToFile)

	hashData := fmt.Sprintf("%s:%s:%s", file.Name(), strconv.FormatInt(file.Size(), 10), file.ModTime())

	hash := utils.Hash([]byte(hashData))

	value, err := env.Database.Read(env.Log, []byte(pathToFile))
	if err != nil {
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
		Name:    file.Name(),
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
