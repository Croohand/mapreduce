package dbase

import (
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

func Open() {
	var err error
	db, err = bolt.Open("bolt.db", 0600, nil)
	if err != nil {
		panic(err)
	}
}

func Close() {
	db.Close()
}

func Set(bucket, key string, value []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		b.Put([]byte(key), value)
		return nil
	})
}

func Get(bucket, key string) (value []byte, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			value = b.Get([]byte(key))
			return nil
		}
		return errors.New(fmt.Sprintf("bucket %s doesn't exist in DB", bucket))
	})
	return
}

func Has(bucket, key string) (bool, error) {
	value, err := Get(bucket, key)
	return value != nil, err
}
