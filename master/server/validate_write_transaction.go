package server

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
	bolt "go.etcd.io/bbolt"
)

func validateWriteTransaction(id, path string, pathInfo fsutil.PathInfo, blocks []byte) (*responses.Answer, error) {
	defer func() {
		transactions.Remove(path)
		for _, block := range pathInfo {
			for _, slave := range block.Slaves {
				http.PostForm(slave+"/Transaction/Remove", url.Values{"TransactionId": {id}})
			}
		}
	}()
	willTry := true
	for _, block := range pathInfo {
		available := 0
		for _, slave := range block.Slaves {
			if !isSlaveAvailable(slave) {
				continue
			}
			resp, err := http.PostForm(slave+"/Block/IsExists", url.Values{"BlockId": {block.Id}, "TransactionId": {id}})
			if err != nil {
				log.Println(err)
				continue
			}
			var bStatus responses.BlockStatus
			if err := httputil.GetJson(resp, &bStatus); err != nil {
				log.Println(err)
				continue
			}
			if bStatus.Exists {
				available += 1
			}
		}
		if available < getMrConfig().MinReplicationFactor {
			willTry = false
			break
		}
	}
	if !willTry {
		return nil, errors.New("can't validate write transaction, relevant slaves are down")
	}
	failed := false
	for _, block := range pathInfo {
		written := 0
		for _, slave := range block.Slaves {
			if !isSlaveAvailable(slave) {
				continue
			}
			resp, err := http.PostForm(slave+"/Block/Validate", url.Values{"BlockId": {block.Id}, "TransactionId": {id}})
			if err != nil {
				log.Println(err)
				continue
			}
			var ans responses.Answer
			if err := httputil.GetJson(resp, &ans); err != nil {
				log.Println(err)
				continue
			}
			if ans.Success {
				written += 1
			}
		}
		if written < getMrConfig().MinReplicationFactor {
			failed = true
			break
		}
	}

	if failed {
		for _, block := range pathInfo {
			for _, slave := range block.Slaves {
				if !isSlaveAvailable(slave) {
					continue
				}
				resp, err := http.PostForm(slave+"/Block/Remove", url.Values{"BlockId": {block.Id}})
				if err != nil {
					log.Println(err)
					continue
				}
				var ans responses.Answer
				if err := httputil.GetJson(resp, &ans); err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}

	if failed {
		return nil, errors.New("failed to validate write transaction, relevant slaves are down")
	}

	err := filesDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Files"))
		if err != nil {
			return err
		}
		b.Put([]byte(path), blocks)
		return nil
	})
	if err != nil {
		return nil, errors.New("failed to write new path to DB")
	}
	return &responses.Answer{true}, nil
}
