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
	"context"
	"errors"
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/filesystem"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type Storage struct {
	lock       sync.Mutex
	root       string
	Collection map[string]map[string]*File
}

type File struct {
	io.Writer
	path string
	dir  string
	file *os.File
}

func (f *File) Write(msg []byte) (int, error) {

	if f.file == nil {
		return 0, errors.New("file descriptor not found")
	}

	msg = append(msg, filesystem.LineSeparator)
	return f.file.Write(msg)
}

func (f *File) Read(ctx context.Context, lines int, follow bool, l chan string) error {

	var (
		seek *tail.SeekInfo
		done bool
	)

	if lines > 0 {
		offset, err := filesystem.LineSeek(lines, f.file)
		if err != nil {
			return err
		}
		seek = new(tail.SeekInfo)
		seek.Offset = offset
		seek.Whence = 1
	}

	t, err := tail.TailFile(f.path, tail.Config{Follow: follow, Location: seek, Logger: tail.DiscardingLogger})
	if err != nil {
		log.Errorf("tail err: %s", err.Error())
		return err
	}

	defer t.Cleanup()

	go func() {
		<-ctx.Done()

		if err := t.Stop(); err != nil {
			log.Errorf("%s:> stop tailing err: %s", logPrefix, err.Error())
		}
		done = true
	}()

	for line := range t.Lines {

		if done {
			break
		}

		l <- line.Text
	}

	fmt.Println("read exit")
	return nil
}

func (f *File) Clear() error {
	if err := f.file.Truncate(0); err != nil {
		return err
	}
	if _, err := f.file.Seek(0, 0); err != nil {
		return err
	}

	return nil
}

func (s Storage) GetStream(kind, selflink string, clear bool) (*File, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.Collection[kind]; !ok {
		s.Collection[kind] = make(map[string]*File)
	}

	if _, ok := s.Collection[kind][selflink]; !ok {
		f, err := newFile(s.root, kind, selflink, clear)
		if err != nil {
			log.Errorf("can not create storage file: %s", err.Error())
			return nil, err
		}
		s.Collection[kind][selflink] = f
	}
	return s.Collection[kind][selflink], nil
}

func NewStorage(root string) *Storage {
	stg := new(Storage)

	stg.root = root
	stg.Collection = make(map[string]map[string]*File)

	return stg
}

func newFile(root, kind, selflink string, clear bool) (*File, error) {
	f := new(File)

	f.dir = filepath.Join(root, kind)
	f.path = filepath.Join(root, kind, selflink)

	if _, err := os.Stat(f.dir); os.IsNotExist(err) {
		err = os.MkdirAll(f.dir, 0755)
		if err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(f.path); err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(f.path)
			if err != nil {
				_ = fmt.Errorf(err.Error())
				return nil, err
			}
		}
	}

	if f.file == nil {
		file, err := os.OpenFile(f.path, os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			_ = fmt.Errorf(err.Error())
			return nil, err
		}
		f.file = file
	}

	if clear {
		if err := f.Clear(); err != nil {
			return nil, err
		}
	}

	return f, nil
}
