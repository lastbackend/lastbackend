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

package mockdb

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/boltdb/bolt"
)

type DB struct {
	db *bolt.DB
}

func Init() (*DB, error) {

	var (
		d = new(DB)
	)

	f, err := ioutil.TempFile("", "")
	if err != nil {
		panic("temp file: " + err.Error())
	}
	path := f.Name()
	defer func() {
		f.Close()
		os.Remove(path)
	}()

	d.db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		panic("open: " + err.Error())
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

		if err := json.Unmarshal(buf, iface); err != nil {
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

		buf, err := json.Marshal(iface)
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
	d.Close()
	return nil
}

func (d *DB) Close() {
	defer os.Remove(d.db.Path())
	d.db.Close()
}
