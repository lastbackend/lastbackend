package db

import (
	"encoding/json"
	"os"
	"path"

	"github.com/boltdb/bolt"
	"github.com/lastbackend/lastbackend/pkg/util/filesystem"
	"github.com/lastbackend/lastbackend/pkg/util/homedir"
)

type DB struct {
	db *bolt.DB
}

func Init() (*DB, error) {

	var (
		err error
		d   = new(DB)
	)

	dir := path.Join(homedir.HomeDir(), string(os.PathSeparator), ".lb")
	err = filesystem.MkDir(dir, 0755)
	if err != nil {
		return nil, err
	}

	d.db, err = bolt.Open(path.Join(dir, string(os.PathSeparator), "lb.db"), 0755, nil)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *DB) Get(fieldname string, iface interface{}) error {

	err := d.db.View(func(tx *bolt.Tx) error {
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

func (d *DB) Set(fieldname string, iface interface{}) error {

	err := d.db.Update(func(tx *bolt.Tx) error {
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

func (d *DB) Clear() error {
	err := os.RemoveAll(path.Join(homedir.HomeDir(), string(os.PathSeparator), ".lb"))
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) Close() error {
	return d.db.Close()
}
