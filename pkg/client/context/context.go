package context

import (
	"github.com/boltdb/bolt"
	"github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/interface/log"
	l "github.com/lastbackend/lastbackend/libs/log"
)

var context Context

func Get() *Context {
	context.Session = new(session)
	return &context
}

func Mock() *Context {

	context.Log = new(l.Log)
	context.Log.Init()
	context.Log.Disabled()

	context.Storage, _ = bolt.Open("/tmp/test.db", 0755, nil)
	context.Session = new(session)

	return &context
}

type Context struct {
	Session *session
	Log     log.ILogger
	HTTP    *http.RawReq
	Storage *bolt.DB
	// Other info for HTTP handlers can be here, like user UUID
}

type session struct {
	Token *string
}

func (s *session) Get() (*string, error) {
	if s.Token != nil {
		return s.Token, nil
	}

	err := context.Storage.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("session"))
		buf := bucket.Get([]byte("token"))
		token := string(buf)
		s.Token = &(token)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.Token, nil
}

func (s *session) Set(token string) error {
	s.Token = &token

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
