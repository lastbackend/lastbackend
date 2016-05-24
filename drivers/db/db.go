package db

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/deployithq/deployit/drivers/interfaces"
)

const bucket = "map"

type Bolt struct {
	DB *bolt.DB
}

func Open(log interfaces.Log, path string) *bolt.DB {

	db, err := bolt.Open(path, 0766, nil)
	if err != nil {
		log.Error(err)
	}

	return db

}

func (b *Bolt) Write(log interfaces.Log, key, value []byte) error {
	log.Debug("Write hash info to database")

	err := b.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Error(err)
			return err
		}

		err = bucket.Put(key, value)
		if err != nil {
			log.Error(err)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (b *Bolt) Read(log interfaces.Log, key []byte) (string, error) {
	log.Debug("Read hash info from database")

	var val []byte

	err := b.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			err := errors.New(fmt.Sprintf("Bucket %s not found", bucket))
			log.Error(err)
			return err
		}

		val = bucket.Get(key)

		return nil
	})

	if err != nil {
		return string(val), err
	}

	return string(val), nil
}
