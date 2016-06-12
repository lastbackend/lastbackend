package utils

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func AppName(path string) string {
	splittedPath := strings.Split(path, "/")
	appName := splittedPath[len(splittedPath)-1]
	return appName
}

// GenerateID - generate random id with 64 length
func GenerateID() string {
	u1 := uuid.NewV4()
	u2 := uuid.NewV4()
	return strings.Replace(fmt.Sprintf("%s%s", u1.String(), u2.String()), "-", "", -1)
}

// TODO: for any exported function comments are necessary.
// For time saving reason just in 2 words put small information
// about what function do and why it is needed for

// Hash - create hash based on provided string
func Hash(data string) string {
	hash := sha1.Sum([]byte(data))

	var hashString string

	for i := 0; i < len(hash); i++ {
		hashString += base64.URLEncoding.EncodeToString(hash[:i])
	}

	return hashString
}

// Ungzip -
func Ungzip(source, target string) error {

	reader, err := os.Open(source)
	if err != nil {
		return err
	}

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}

	target = filepath.Join(target, archive.Name)

	writer, err := os.Create(target)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, archive)

	reader.Close()
	archive.Close()
	writer.Close()

	return err
}

func Update(source, target, update string, excludes []string) error {

	src_f, err := os.Open(source)
	if err != nil {
		return err
	}

	src_rd := tar.NewReader(src_f)

	target_f, err := os.Create(target)
	if err != nil {
		return err
	}

	target_wr := tar.NewWriter(target_f)
	updated := make(map[string]bool)

	_, err = os.Stat(update)

	// apply the update if there is file
	if update != `` && !os.IsNotExist(err) {
		update_f, err := os.Open(update)
		if err != nil {
			return err
		}

		update_rd := tar.NewReader(update_f)

		// write new files
		for {
			header, err := update_rd.Next()
			if err == io.EOF {
				target_wr.Flush()
				break
			} else if err != nil {
				return err
			} else {
				if _, ok := updated[header.Name]; !ok {
					updated[header.Name] = true
				}

				// write the header to the tarball archive
				if err := target_wr.WriteHeader(header); err != nil {
					return err
				}

				// replicate the file/dir to the tarball
				if _, err := io.Copy(target_wr, update_rd); err != nil {
					return err
				}
			}
		}

		update_f.Close()
	}

	// merge new files with current layer and exclude remove files
	for {
		header, err := src_rd.Next()

		if err == io.EOF {
			target_wr.Flush()
			break
		} else if err != nil {
			return err
		} else {
			if _, ok := updated[header.Name]; !ok {

				path := header.Name

				if header.Typeflag == tar.TypeDir {
					path = trimSuffix(path, "/")
				}

				if !exists(excludes, path) {
					// write the header to the tarball archive
					if err := target_wr.WriteHeader(header); err != nil {
						return err
					}

					// replicate the file/dir to the tarball
					if _, err := io.Copy(target_wr, src_rd); err != nil {
						return err
					}
				}
			}
		}
	}

	src_f.Close()
	target_f.Close()
	target_wr.Close()

	return nil
}

func ReadConfig(path string, i interface{}) error {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return nil
	}

	reader, err := os.Open(path)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if header.Name == `deployit.yaml` {

			data, err := ioutil.ReadAll(tarReader)
			if err != nil {
				return err
			}

			if err := yaml.Unmarshal(data, i); err != nil {
				return err
			}

			break
		}
	}

	reader.Close()

	return nil
}

func CreateDirs(paths []string) error {

	fileMode := os.FileMode(666)

	for _, path := range paths {
		if err := os.MkdirAll(path, fileMode); err != nil {
			return err
		}
	}

	return nil
}

func RemoveDirs(paths []string) error {

	for _, path := range paths {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}

// Returns `true` if the target string t is in the
// slice.
func exists(vs []string, t string) bool {
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

func GetStringParamFromURL(param string, r *http.Request) string {
	params := mux.Vars(r)
	value := params[param]
	return value
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		rand.Seed(int64(time.Now().Nanosecond()))
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func FileLine() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
