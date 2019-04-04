package server

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
)

func isSlaveAvailable(addr string) bool {
	resp, err := http.Get(addr + "/IsAlive")
	if err == nil {
		var status responses.SlaveStatus
		if httputil.GetJson(resp, &status) == nil && status.Alive {
			return true
		}
	}
	return false
}

func getAvailableSlaves(lim int) ([]string, error) {
	slaves := make([]string, 0)
	for _, i := range rand.Perm(len(Config.SlaveAddrs)) {
		addr := Config.SlaveAddrs[i]
		if isSlaveAvailable(addr) {
			slaves = append(slaves, addr)
			if len(slaves) == lim {
				return slaves, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("not enough slaves, asked %d got %d", lim, len(slaves)))
}
