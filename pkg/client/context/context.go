package context

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/interface/log"
	l "github.com/lastbackend/lastbackend/libs/log"
	f "github.com/lastbackend/lastbackend/utils"
	"os"
)

var context Context

func Get() *Context {
	return &context
}

func Mock() *Context {
	context.mock = true
	context.Log = new(l.Log)
	context.Log.Init()
	context.Log.Disabled()

	return &context
}

type Context struct {
	Log     log.ILogger
	HTTP    *http.RawReq
	Storage ILocalStorage
	mock    bool
	// Other info for HTTP handlers can be here, like user UUID
}

type ILocalStorage interface {
	Get(string, interface{}) error
	Set(string, interface{}) error
	Clear() error
	Init() error
}

type LocalStorage struct {
	db *bolt.DB
}

func (s *LocalStorage) Init() error {

	dir := f.GetHomeDir() + "/.lb"

	f.MkDir(dir, 0755)

	var err error
	s.db, err = bolt.Open(dir+"/lb.db", 0755, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *LocalStorage) Get(fieldname string, iface interface{}) error {

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("storage"))
		if bucket == nil {
			return nil
		}

		buf := bucket.Get([]byte(fieldname))
		if buf == nil {
			return nil
		}

		err := json.Unmarshal(buf, iface)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *LocalStorage) Set(fieldname string, iface interface{}) error {

	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("storage"))
		if err != nil {
			return err
		}

		ifacebyte, err := json.Marshal(&iface)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(fieldname), ifacebyte)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *LocalStorage) Clear() error {

	fmt.Println("lal")

	err := os.RemoveAll(f.GetHomeDir() + "/.lb")
	if err != nil {
		return err
	}

	return nil
}
