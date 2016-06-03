package utils

import (
	"os"
	"io"
	"archive/tar"
)

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
