package utils

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha1"
	"encoding/base64"
	"github.com/deployithq/deployit/drivers/interfaces"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Hash(data string) string {
	hash := sha1.Sum([]byte(data))

	var hashString string

	for i := 0; i < len(hash); i++ {
		hashString += base64.URLEncoding.EncodeToString(hash[:i])
	}

	return hashString
}

func Untar(log interfaces.ILog, filename string) error {
	file, err := os.Open(filename)

	if err != nil {
		log.Error(err)
		return err
	}

	defer file.Close()

	var fileReader io.ReadCloser = file

	// just in case we are reading a tar.gz file, add a filter to handle gzipped file
	if strings.HasSuffix(filename, ".gz") {
		if fileReader, err = gzip.NewReader(file); err != nil {
			log.Error(err)
			return err
		}

		defer fileReader.Close()
	}

	tarBallReader := tar.NewReader(fileReader)

	// Extracting tarred files

	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Error(err)
		}

		// get the individual filename and extract to the current directory
		filename := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory
			log.Debug("Creating directory :", filename)
			err = os.MkdirAll(filename, os.FileMode(header.Mode)) // or use 0755 if you prefer

			if err != nil {
				log.Error(err)
				return err
			}

		case tar.TypeReg:
			// handle normal file
			log.Debug("Untarring :", filename)
			writer, err := os.Create(filename)

			if err != nil {
				log.Error(err)
				return err
			}

			io.Copy(writer, tarBallReader)

			err = os.Chmod(filename, os.FileMode(header.Mode))

			if err != nil {
				log.Error(err)
				return err
			}

			writer.Close()
		default:
			log.Debugf("Unable to untar type : %c in file %s", header.Typeflag, filename)
		}
	}

	return nil
}

func Ungzip(log interfaces.ILog, source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	target = filepath.Join(target, archive.Name)

	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}
