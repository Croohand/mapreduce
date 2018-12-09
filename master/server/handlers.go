package server

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/Croohand/mapreduce/common/httputil"
)

const maxTimeout = 100 * time.Millisecond

var fastClient = http.Client{Timeout: time.Duration(maxTimeout)}

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
	const maxAliveCheckRetries = 5
	for i := 0; i < maxAliveCheckRetries; i++ {
		resp, err := fastClient.Get(addr + "/IsAlive")
		if err == nil {
			var alive struct{ Alive bool }
			if httputil.GetJson(resp, &alive) == nil && alive.Alive {
				return true
			}
		}
		time.Sleep(time.Duration(maxTimeout))
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
