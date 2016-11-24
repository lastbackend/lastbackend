package context

import (
	"github.com/boltdb/bolt"
	"github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/interface/log"
	l "github.com/lastbackend/lastbackend/libs/log"
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
	context.Storage = new(bolt.DB)

	return &context
}

type Context struct {
	Session session
	Log     log.ILogger
	HTTP    *http.RawReq
	Storage *bolt.DB
	mock    bool
	// Other info for HTTP handlers can be here, like user UUID
}

type session struct {
	token *string
}

func (s *session) Get() (*string, error) {

	if s.token != nil || context.mock {
		return s.token, nil
	}

	err := context.Storage.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("session"))
		if bucket == nil {
			return nil
		}

		buf := bucket.Get([]byte("token"))
		token := string(buf)
		s.token = &(token)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.token, nil
}

func (s *session) Set(token string) error {

	s.token = &token

	if context.mock {
		return nil
	}

	err := context.Storage.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("session"))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte("token"), []byte(token))
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

func (s *session) Clear() error {

	s.token = nil

	if context.mock {
		return nil
	}

	err := context.Storage.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("session"))
		if bucket == nil {
			return nil
		}

		err := bucket.Delete([]byte("token"))
		if err != nil {
			return err
		}

		err = bucket.DeleteBucket([]byte("session"))
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
