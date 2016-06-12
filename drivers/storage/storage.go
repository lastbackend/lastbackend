package storage

import (
	"errors"
	"github.com/boltdb/bolt"
	"github.com/deployithq/deployit/drivers/interfaces"
)

const base = "map"

type Storage struct {
	Driver *bolt.DB
}

func Open(path string) (*Storage, error) {

	db, err := bolt.Open(path, 0766, nil)
	if err != nil {
		return new(Storage), err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(base))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return new(Storage), err
	}

	return &Storage{db}, nil

}

func (b *Storage) Write(key, value string) error {

	err := b.Driver.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(base))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(key), []byte(value))
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

func (b *Storage) Read(key string) (string, error) {

	var val string

	err := b.Driver.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(base))
		if bucket == nil {
			err := errors.New("BUCKET_NOT_FOUND")
			return err
		}

		val = string(bucket.Get([]byte(key)))

		return nil
	})

	if err != nil {
		return val, err
	}

	return val, nil
}

func (b *Storage) Delete(key string) error {

	err := b.Driver.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(base))
		if bucket == nil {
			return interfaces.ErrBucketNotFound
		}

		err := bucket.Delete([]byte(key))
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

func (b *Storage) ListAllFiles() (map[string]string, error) {

	files := make(map[string]string)

	err := b.Driver.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(base))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if k != nil {
				files[string(k)] = string(v)
			} else {
				break
			}
		}

		return nil
	})

	if err != nil {
		return files, err
	}

	return files, nil
}
