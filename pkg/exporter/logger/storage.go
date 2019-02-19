//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package logger

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/lastbackend/lastbackend/pkg/log"
	"io"
	"os"
	"path/filepath"
)

type Storage struct {
	root       string
	Collection map[string]map[string]*File
}

type File struct {
	io.Writer
	path string
	dir  string
	file *os.File
}

func (f *File) Write(msg string) error {

	if f.file == nil {
		return errors.New("file descriptor not found")
	}

	f.file.Write([]byte(msg))
	f.file.Write([]byte("\n"))

	return nil
}

func (f *File) ReadLines(count int, clear bool) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(f.file)
	var i = 0
	for scanner.Scan() {
		if count > 0 {
			if i >= count {
				break
			}
			i++
		}
		fmt.Println("read", scanner.Text())
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func (f *File) Tail(count int, clear bool, buffer *bytes.Buffer) error {

	t, err := tail.TailFile(f.path, tail.Config{Follow: true})
	if err != nil {
		log.Errorf("tail err: %s", err.Error())
		return err
	}
	for line := range t.Lines {
		fmt.Println(">tail:", line.Text)
		buffer.Write([]byte(line.Text))
	}
	return nil
}

func (s Storage) GetStream(kind, selflink string) *File {

	if _, ok := s.Collection[kind]; !ok {
		s.Collection[kind] = make(map[string]*File)
	}

	if _, ok := s.Collection[kind][selflink]; !ok {
		fmt.Println("create new fie")
		f, err := newFile(s.root, kind, selflink)
		if err != nil {
			log.Errorf("can not create storage file: %s", err.Error())
			return nil
		}
		s.Collection[kind][selflink] = f
	}
	return s.Collection[kind][selflink]
}

func NewStorage(root string) *Storage {
	stg := new(Storage)

	stg.root = root
	stg.Collection = make(map[string]map[string]*File)

	return stg
}

func newFile(root, kind, selflink string) (*File, error) {
	f := new(File)

	f.dir = filepath.Join(root, kind)
	f.path = filepath.Join(root, kind, selflink)

	fmt.Println("dir:", f.dir, " path:", f.path)

	if _, err := os.Stat(f.dir); os.IsNotExist(err) {
		err = os.MkdirAll(f.dir, 0755)
		if err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(f.path); os.IsNotExist(err) {
		fmt.Println("create new file:", f.path)
		file, err := os.Create(f.path)
		if err != nil {
			fmt.Errorf(err.Error())
			return nil, err
		}
		fmt.Println("file is created:", f.path)
		f.file = file
	}

	if f.file == nil {
		file, err := os.OpenFile(f.path, os.O_APPEND, os.ModeAppend)
		if err != nil {
			fmt.Errorf(err.Error())
			return nil, err
		}
		fmt.Println("open file:", f.path)
		f.file = file
	}

	return f, nil
}
