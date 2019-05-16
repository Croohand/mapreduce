package dbase

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	"github.com/Croohand/mapreduce/common/wrrors"
	bolt "go.etcd.io/bbolt"
)

const (
	Txs     = "Transactions"
	PathTxs = "PathTransactions"
	Files   = "Files"
	Blocks  = "Blocks"
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
		return b.Put([]byte(key), value)
	})
}

func Get(bucket, key string) (value []byte, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			value = b.Get([]byte(key))
			return nil
		}
		return nil
	})
	return
}

func Has(bucket, key string) (bool, error) {
	value, err := Get(bucket, key)
	return value != nil, err
}

func Del(bucket, key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return b.Delete([]byte(key))
	})
}

func Marshal(v interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func Unmarshal(data []byte, v interface{}) error {
	b := bytes.NewBuffer(data)
	return gob.NewDecoder(b).Decode(v)
}

func GetObject(bucket, key string, res interface{}) error {
	data, err := Get(bucket, key)
	if err != nil {
		return err
	}
	if data == nil {
		return wrrors.New("GetObject").WrapS(fmt.Sprintf("No object with key %s in bucket %s", key, bucket))
	}
	return json.Unmarshal(data, res)
}

func SetObject(bucket, key string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	err = Set(bucket, key, data)
	if err != nil {
		return err
	}
	return nil
}

func GetKeys(bucket, prefix string) ([]string, error) {
	var result []string
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, _ := c.Seek([]byte(prefix)); k != nil && bytes.HasPrefix(k, []byte(prefix)); k, _ = c.Next() {
			result = append(result, string(k))
		}
		return nil
	})
	return result, nil
}
