package context

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/interface/localdb"
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
	context.Storage = new(LocalStorage)

	return &context
}

type Context struct {
	Log     log.ILogger
	HTTP    *http.RawReq
	Storage localdb.ILocalStorage
	mock    bool
	// Other info for HTTP handlers can be here, like user UUID
}

// TODO: It is necessary to move to libs
type LocalStorage struct {
	db *bolt.DB
}

func (s *LocalStorage) Init() error {

	var err error

	dir := f.GetHomeDir() + "/.lb"
	err = f.MkDir(dir, 0755)
	if err != nil {
		return err
	}

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

		buf, err := json.Marshal(&iface)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(fieldname), buf)
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
	err := os.RemoveAll(f.GetHomeDir() + "/.lb")
	if err != nil {
		return err
	}

	return nil
}

func (s *LocalStorage) Close() error {
	return s.db.Close()
}
