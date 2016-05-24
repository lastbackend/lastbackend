package db

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const bucket = "map"

type Bolt struct {
	DB *bolt.DB
}

func Open() *bolt.DB {

	db, err := bolt.Open("/home/nate/foo/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}

	return db

}

func (b *Bolt) Write(key, value []byte) {

	err := b.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}

		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (b *Bolt) Read(key []byte) {

	err := b.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}

		val := bucket.Get(key)
		fmt.Println(string(val))

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
