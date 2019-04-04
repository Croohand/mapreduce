package server

import (
	"errors"

	"github.com/Croohand/mapreduce/common/responses"
	bolt "go.etcd.io/bbolt"
)

func isFileExists(path string) (r *responses.FileStatus, err error) {
	r = &responses.FileStatus{}
	err = filesDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Files"))
		if b != nil {
			r.Exists = b.Get([]byte(path)) != nil
			return nil
		}
		return errors.New("bucket Files doesn't exist in DB")
	})
	return
}
