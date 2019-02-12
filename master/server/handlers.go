package server

import (
	"math/rand"
	"net/http"

	"github.com/Croohand/mapreduce/common/httputil"
	bolt "go.etcd.io/bbolt"
)

func isAlive(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJson(w, struct {
		Alive bool
		Type  string
	}{true, "master"})
}

func getMRConfig(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJson(w, mrConfig)
}

func isSlaveAvailable(addr string) bool {
	resp, err := http.Get(addr + "/IsAlive")
	if err == nil {
		var alive struct{ Alive bool }
		if httputil.GetJson(resp, &alive) == nil && alive.Alive {
			return true
		}
	}
	return false
}

func getAvailableSlavesInner(lim int) (slaves []string) {
	for _, i := range rand.Perm(len(Config.SlaveAddrs)) {
		addr := Config.SlaveAddrs[i]
		if isSlaveAvailable(addr) {
			slaves = append(slaves, addr)
			if len(slaves) == lim {
				break
			}
		}
	}
	return
}

func getAvailableSlaves(w http.ResponseWriter, r *http.Request) {
	slaves := getAvailableSlavesInner(mrConfig.ReplicationFactor)
	if len(slaves) < mrConfig.ReplicationFactor {
		http.Error(w, "not enough available slaves for write", http.StatusServiceUnavailable)
		return
	}
	httputil.WriteJson(w, struct{ Slaves []string }{slaves})
}

func isFileExists(w http.ResponseWriter, r *http.Request) {
	path := r.PostFormValue("Path")
	if path == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	exists, failed := false, false
	err := filesDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Files"))
		if b != nil {
			exists = b.Get([]byte(path)) != nil
		} else {
			http.Error(w, "bucket Files doesn't exist in DB", http.StatusInternalServerError)
			failed = true
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		failed = true
	}
	if !failed {
		httputil.WriteJson(w, struct{ Exists bool }{exists})
	}
}
