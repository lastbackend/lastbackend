package handlers

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/deployithq/deployit/drivers/interfaces"
	print_ "github.com/deployithq/deployit/drivers/print"
	"github.com/deployithq/deployit/drivers/storage"
	"github.com/deployithq/deployit/errors"
	"github.com/deployithq/deployit/utils"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ItCommand struct {
	Host struct {
		Name string
		URL  string
	}
	Paths struct {
		Root    string
		Storage string
		Archive string
	}
	Tag     string
	Storage interfaces.IStorage
	Print   interfaces.IPrint
}

func (c *ItCommand) Run(args []string) int {

	var err error

	// Initializaing printing module
	c.Print = print_.Init()

	// Adding path module
	c.Paths.Root, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		c.Print.Error(err)
		return 1
	}

	// Adding storage module
	c.Paths.Storage = fmt.Sprintf("%s/.dit", c.Paths.Root)

	c.Paths.Archive = fmt.Sprintf("%s/.dit/%s", c.Paths.Root, "tar.gz")

	// Creating storage directory
	err = os.Mkdir(c.Paths.Storage, 0766)
	if err != nil && os.IsNotExist(err) {
		c.Print.Error(err)
		return 1
	}

	// Initializaing storage module
	c.Storage, err = storage.Open(fmt.Sprintf("%s/%s_map", c.Paths.Storage, c.Host))
	if err != nil {
		c.Print.Error(err)
		return 1
	}

	var debug bool

	// Creating flags set
	cmdFlags := flag.NewFlagSet("it", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		c.Print.WhiteInfo(c.Help())
	}

	cmdFlags.BoolVar(&debug, "debug", false, "Enables debug mode")
	if Debug == false {
		if os.Getenv("DEPLOYIT_DEBUG") != "" {
			Debug = true
		}
	}

	cmdFlags.StringVar(&c.Host.URL, "host-url", "https://api.deployit.io", "URL of host, where daemon is running")
	if Debug == false {
		if os.Getenv("DEPLOYIT_HOST_URL") != "" {
			Host = os.Getenv("DEPLOYIT_HOST_URL")
		}
	}

	// Parsing flags
	if err := cmdFlags.Parse(args); err != nil {
		c.Print.WhiteInfo(c.Help())
		return 1
	}

	// Setting debug mode to printer if debug == true
	c.Print.SetDebug(debug)

	// Parsing url
	u, err := url.Parse(c.Host.URL)
	if err != nil {
		c.Print.Error(err)
		return 1
	}

	c.Host.Name = u.Host

	c.Tag = "latest"

	// Main function
	if err := DeployIt(c); err != nil {
		return 1
	}

	return 0

}

func (c *ItCommand) Help() string {
	return ""
}

func (c *ItCommand) Synopsis() string {
	return ""
}

// Main deploy it handler
// - Archive all files which is in folder
// - Send it to server

func DeployIt(c *ItCommand) error {

	appInfo := new(AppInfo)
	err := appInfo.Read(c.Paths.Root, c.Host.Name)
	if err != nil {
		c.Print.Error(err)
		return err
	}

	if appInfo.Name == "" {

		appInfo.Name = utils.AppName(c.Paths.Root)
		appInfo.Tag = c.Tag

		err = appCreate(c, appInfo.Name, appInfo.Tag)
		if err != nil {
			return err
		}

		c.Print.Infof("Creating app: %s", appInfo.Name)
	} else {
		c.Print.Infof("Updating app: %s", appInfo.Name)
	}

	fw, gw, tw, err := utils.CreateTarGz(c.Paths.Archive)
	if err != nil {
		c.Print.Error(err)
		return err
	}

	// Deleting archive after function ends
	defer os.Remove(c.Paths.Archive)

	// Listing all files from database to know what files were deleted from previous run
	storedFiles, err := c.Storage.ListAllFiles()
	if err != nil {
		c.Print.Error(err)
		return err
	}

	// TODO Include deleted folders to deletedFiles like "nginx/"

	excludePatterns, err := utils.LoadDockerPatterns(c.Paths.Root)
	if err != nil {
		c.Print.Error(err)
		return err
	}

	excludePatterns = append(excludePatterns, ".gitignore", ".dit", ".git")

	c.Print.Info("Packing files")
	storedFiles, err = packFiles(c, tw, c.Paths.Root, storedFiles, excludePatterns)
	if err != nil {
		c.Print.Error(err)
		return err
	}

	fw.Close()
	gw.Close()
	tw.Close()

	deletedFiles := []string{}

	for key, _ := range storedFiles {
		c.Print.Debug("Deleting: ", key)
		err = c.Storage.Delete(key)
		if err != nil {
			c.Print.Error(err)
			return err
		}

		deletedFiles = append(deletedFiles, key)
	}

	bodyBuffer := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuffer)

	// Adding deleted files to request
	if len(deletedFiles) > 0 {
		delFiles, err := json.Marshal(deletedFiles)
		if err != nil {
			c.Print.Error(err)
			return err
		}

		bodyWriter.WriteField("deleted", string(delFiles))
	}

	archiveInfo, err := os.Stat(c.Paths.Archive)
	if err != nil {
		c.Print.Error(err)
		return err
	}

	// If archive size is 32 it means that it is empty and we don't need to send it
	if archiveInfo.Size() != 32 {
		fh, err := os.Open(c.Paths.Archive)
		if err != nil {
			c.Print.Error(err)
			return err
		}

		fileWriter, err := bodyWriter.CreateFormFile("file", "tar.gz")
		if err != nil {
			c.Print.Error(err)
			return err
		}

		_, err = io.Copy(fileWriter, fh)
		if err != nil {
			c.Print.Error(err)
			return err
		}

		fh.Close()
	}

	bodyWriter.Close()

	// TODO If error in response: rollback hash table

	c.Print.Info("Uploading sources")

	res, err := utils.ExecuteHTTPRequest(fmt.Sprintf("%s/app/%s/deploy", c.Host.URL, appInfo.Name), "POST",
		bodyWriter.FormDataContentType(), bodyBuffer)

	if err != nil {
		c.Print.Error(err)
		return err
	}

	c.Print.Info("Logs: ")
	utils.StreamHttpResponse(res)

	err = appInfo.Write(c.Paths.Root, c.Host.Name, appInfo.Name, appInfo.Tag, appInfo.URL)
	if err != nil {
		c.Print.Error(err)
		return err
	}

	res.Body.Close()

	c.Print.Info("Done")

	return nil
}

func packFiles(c *ItCommand, tw *tar.Writer, filesPath string, storedFiles map[string]string, excludePatterns []string) (map[string]string, error) {

	// Opening directory with files
	dir, err := os.Open(filesPath)

	if err != nil {
		c.Print.Error(err)
		return storedFiles, err
	}

	// Reading all files
	files, err := dir.Readdir(0)
	if err != nil {
		c.Print.Error(err)
		return storedFiles, err
	}

	for _, file := range files {

		currentFilePath := fmt.Sprintf("%s/%s", filesPath, file.Name())

		// Creating path, which will be inside of archive
		relativePath := strings.Replace(currentFilePath, c.Paths.Root, "", 1)[1:]

		// Ignoring files which is not needed for build to make archive smaller
		// TODO: create base .ditignore file on first application creation
		matches, err := utils.Matches(relativePath, excludePatterns)
		if err != nil {
			c.Print.Error(err)
			return storedFiles, err
		}

		if matches {
			continue
		}

		// If it was directory - calling this function again
		// In other case adding file to archive
		if file.IsDir() {
			storedFiles, err = packFiles(c, tw, currentFilePath, storedFiles, excludePatterns)
			if err != nil {
				return storedFiles, err
			}
			continue
		}

		// Creating hash
		hash := utils.Hash(fmt.Sprintf("%s:%s:%s", file.Name(), strconv.FormatInt(file.Size(), 10), file.ModTime()))

		// If file has not changed - don't add it to archive
		if storedFiles[relativePath] == hash {
			delete(storedFiles, relativePath)
			continue
		}

		delete(storedFiles, relativePath)

		// If hashes are not equal - add file to archive
		c.Print.Debug("Packing file: ", currentFilePath)

		// Adding file to hash table
		err = c.Storage.Write(relativePath, hash)
		if err != nil {
			c.Print.Error(err)
			return storedFiles, err
		}

		// Adding file to archive
		err = utils.AddFileToArchive(tw, file, currentFilePath, relativePath)
		if err != nil {
			c.Print.Error(err)
			return storedFiles, err
		}

	}

	dir.Close()

	return storedFiles, err

}

func appCreate(c *ItCommand, name, tag string) error {

	request := struct {
		Name string `json:"name"`
		Tag  string `json:"tag"`
	}{name, tag}

	var buf io.ReadWriter
	buf = new(bytes.Buffer)

	err := json.NewEncoder(buf).Encode(request)
	if err != nil {
		return err
	}

	res, err := utils.ExecuteHTTPRequest(fmt.Sprintf("%s/app", c.Host.URL), "PUT", "application/json; charset=utf-8", buf)
	if err != nil {
		c.Print.Error(err)
		return err
	}

	if res.StatusCode != 200 {
		err = errors.ParseError(res)
		c.Print.Error(err)
		return err
	}

	return nil
}
