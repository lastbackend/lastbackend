package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
)

func CreateTarGz(pathToArchive string) (*os.File, *gzip.Writer, *tar.Writer, error) {

	gw := new(gzip.Writer)
	tw := new(tar.Writer)

	fw, err := os.Create(pathToArchive)
	if err != nil {
		return fw, gw, tw, err
	}

	gw = gzip.NewWriter(fw)
	tw = tar.NewWriter(gw)

	return fw, gw, tw, nil

}

func AddFileToArchive(tw *tar.Writer, file os.FileInfo, currentFilePath, relativePath string) error {

	fr, err := os.Open(currentFilePath)
	if err != nil {
		return err
	}

	h := &tar.Header{
		Name:    relativePath,
		Size:    file.Size(),
		Mode:    int64(file.Mode()),
		ModTime: file.ModTime(),
	}

	err = tw.WriteHeader(h)
	if err != nil {
		return err
	}

	_, err = io.Copy(tw, fr)
	if err != nil {
		return err
	}

	fr.Close()

	return nil
}
