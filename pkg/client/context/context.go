package context

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/interface/log"
	l "github.com/lastbackend/lastbackend/libs/log"
	f "github.com/lastbackend/lastbackend/utils"
	"os"
	"fmt"
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
	context.Storage.db = new(bolt.DB)

	return &context
}

type Context struct {
	Log     log.ILogger
	HTTP    *http.RawReq
	Storage storage
	mock    bool
	// Other info for HTTP handlers can be here, like user UUID
}

type storage struct {
	db *bolt.DB
}

func (s *storage) Init() error {

	dir := f.GetHomeDir() + "/.lb"

	f.MkDir(dir, 0755)

	var err error
	s.db, err = bolt.Open(dir+"/lb.db", 0755, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) Get(fieldname string, iface interface{}) error {

	if context.mock {
		return nil
	}

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

func (s *storage) Set(fieldname string, iface interface{}) error {

	fmt.Println("fn: ", fieldname, "if: ", iface)

	if context.mock {
		return nil
	}

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

func (s *storage) Clear() error {
	err := os.RemoveAll(f.GetHomeDir() + "/.lb")
	if err != nil {
		return err
	}

	return nil
}
