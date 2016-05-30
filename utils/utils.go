package utils

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/satori/go.uuid"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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

func CreateLayer(source, target string) error {

	source_f, err := os.Open(source)
	if err != nil {
		return err
	}

	source_rd := tar.NewReader(source_f)

	target_f, err := os.Create(target)
	if err != nil {
		return err
	}

	target_wr := tar.NewWriter(target_f)

	for {
		header, err := source_rd.Next()
		if err == io.EOF {
			target_wr.Flush()
			break
		} else if err != nil {
			return err
		} else if header.Size > 1e6 {
			replicate(header, target_wr, source_rd)
		}
	}

	source_f.Close()
	target_f.Close()
	target_wr.Close()

	return nil
}

func ModifyLayer(source, update, target string, excludes []string) error {

	src_f, err := os.Open(source)
	if err != nil {
		return err
	}

	src_rd := tar.NewReader(src_f)

	update_f, err := os.Open(update)
	if err != nil {
		return err
	}

	update_rd := tar.NewReader(update_f)

	target_f, err := os.Create(target)
	if err != nil {
		return err
	}

	target_wr := tar.NewWriter(target_f)

	excluded := make(map[string]bool)

	for {
		header, err := update_rd.Next()

		if err == io.EOF {
			target_wr.Flush()
			break
		} else if err != nil {
			return err
		} else if header.Size > 1e6 {

			path := header.Name

			if header.Typeflag == tar.TypeDir {
				path = trimSuffix(path, "/")
			}

			if _, ok := excluded[header.Name]; !ok {
				excluded[header.Name] = true
			}

			if !exists(excludes, path) {
				replicate(header, target_wr, update_rd)
			}
		}

		return nil
	}

	for {
		header, err := update_rd.Next()

		if err == io.EOF {
			target_wr.Flush()
			break
		} else if err != nil {
			return err
		} else if header.Size > 1e6 {
			if _, ok := excluded[header.Name]; !ok {
				replicate(header, target_wr, src_rd)
			}
		}
	}

	src_f.Close()
	update_f.Close()
	target_f.Close()
	target_wr.Close()

	return nil
}

func replicate(hdr *tar.Header, wr *tar.Writer, src *tar.Reader) error {
	// write the header to the tarball archive
	if err := wr.WriteHeader(hdr); err != nil {
		return err
	}

	// replicate the file/dir to the tarball
	if _, err := io.Copy(wr, src); err != nil {
		return err
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
