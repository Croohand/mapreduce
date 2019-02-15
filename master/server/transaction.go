package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	bolt "go.etcd.io/bbolt"
)

const transactionWaitTime = 60

func updateTransaction(w http.ResponseWriter, r *http.Request) {
	id, path := r.PostFormValue("TransactionId"), r.PostFormValue("Path")
	if id == "" || path == "" {
		http.Error(w, "empty id or path in transaction check", http.StatusBadRequest)
		return
	}
	if !fsutil.ValidateTransactionId(id) {
		http.Error(w, "invalid transaction id "+id, http.StatusBadRequest)
		return
	}
	v, ok := transactions.Get(path)
	if ok {
		tx, ok := v.(Transaction)
		if !ok {
			transactions.Remove(path)
		} else if time.Since(tx.LastUpdate).Seconds() > transactionWaitTime {
			transactions.Remove(path)
			if tx.Id == id {
				http.Error(w, "Transaction is too old", http.StatusInternalServerError)
				return
			}
		} else if tx.Id != id {
			httputil.WriteJson(w, struct{ Success bool }{false})
			return
		}
	}
	transactions.Set(path, Transaction{Id: id, LastUpdate: time.Now()})
	httputil.WriteJson(w, struct{ Success bool }{true})
}

func isAliveTransactionInner(id, path string) bool {
	v, ok := transactions.Get(path)
	if ok {
		tx, ok := v.(Transaction)
		if !ok || time.Since(tx.LastUpdate).Seconds() > transactionWaitTime {
			transactions.Remove(path)
			return false
		}
		return tx.Id == id
	}
	return false
}

func isAliveTransaction(w http.ResponseWriter, r *http.Request) {
	id, path := r.PostFormValue("TransactionId"), r.PostFormValue("Path")
	if id == "" || path == "" {
		http.Error(w, "empty id or path in transaction check", http.StatusBadRequest)
		return
	}
	if !fsutil.ValidateTransactionId(id) {
		http.Error(w, "invalid transaction id "+id, http.StatusBadRequest)
		return
	}
	httputil.WriteJson(w, struct{ Alive bool }{isAliveTransactionInner(id, path)})
}

func validateWriteTransaction(w http.ResponseWriter, r *http.Request) {
	id, path, blocks := r.PostFormValue("TransactionId"), r.PostFormValue("Path"), r.PostFormValue("PathInfo")
	if !fsutil.ValidateTransactionId(id) {
		http.Error(w, "invalid transaction id "+id, http.StatusBadRequest)
		return
	}
	if !fsutil.ValidateFilePath(path) {
		http.Error(w, "invalid file path "+path, http.StatusBadRequest)
		return
	}
	var pathInfo fsutil.PathInfo
	if err := json.Unmarshal([]byte(blocks), &pathInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func() {
		transactions.Remove(path)
		for _, block := range pathInfo {
			for _, slave := range block.Slaves {
				http.PostForm(slave+"/Transaction/Remove", url.Values{"TransactionId": []string{id}})
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
			resp, err := http.PostForm(slave+"/Block/IsExists", url.Values{"BlockId": []string{block.Id}, "TransactionId": []string{id}})
			if err != nil {
				log.Println("IsExists block error: " + err.Error())
				continue
			}
			var exists struct{ Exists bool }
			if err := httputil.GetJson(resp, &exists); err != nil {
				log.Println("IsExists block error: " + err.Error())
				continue
			}
			if exists.Exists {
				available += 1
			}
		}
		if available < mrConfig.MinReplicationFactor {
			willTry = false
			break
		}
	}
	if !willTry {
		http.Error(w, "can't validate write transaction, relevant slaves are down", http.StatusInternalServerError)
		return
	}
	failed := false
	for _, block := range pathInfo {
		written := 0
		for _, slave := range block.Slaves {
			if !isSlaveAvailable(slave) {
				continue
			}
			resp, err := http.PostForm(slave+"/Block/Validate", url.Values{"BlockId": []string{block.Id}, "TransactionId": []string{id}})
			if err != nil {
				log.Println("validate block error: " + err.Error())
				continue
			}
			var success struct{ Success bool }
			if err := httputil.GetJson(resp, &success); err != nil {
				log.Println("validate block error: " + err.Error())
				continue
			}
			if success.Success {
				written += 1
			}
		}
		if written < mrConfig.MinReplicationFactor {
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
				resp, err := http.PostForm(slave+"/Block/Remove", url.Values{"BlockId": []string{block.Id}})
				if err != nil {
					log.Println("remove block error: " + err.Error())
					continue
				}
				var success struct{ Success bool }
				if err := httputil.GetJson(resp, &success); err != nil {
					log.Println("remove block error: " + err.Error())
					continue
				}
			}
		}
	}

	if failed {
		http.Error(w, "failed to validate write transaction, relevant slaves are down", http.StatusInternalServerError)
		return
	}

	err := filesDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Files"))
		if err != nil {
			return err
		}
		b.Put([]byte(path), []byte(blocks))
		return nil
	})
	if err != nil {
		http.Error(w, "failed to write new path to DB", http.StatusInternalServerError)
	}
	httputil.WriteJson(w, struct{ Success bool }{true})
}
