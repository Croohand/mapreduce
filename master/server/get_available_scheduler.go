package server

import (
	"errors"
	"math/rand"

	"github.com/Croohand/mapreduce/common/httputil"
)

func getAvailableScheduler() (string, error) {
	for _, i := range rand.Perm(len(Config.SchedulerAddrs)) {
		addr := Config.SchedulerAddrs[i]
		if httputil.IsSlaveAvailable(addr) {
			return addr, nil
		}
	}
	return "", errors.New("No schedulers available")
}
