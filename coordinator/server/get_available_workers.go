package server

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/Croohand/mapreduce/common/httputil"
)

func getAvailableWorkers(lim int) ([]string, error) {
	workers := make([]string, 0)
	for _, i := range rand.Perm(len(Config.WorkerAddrs)) {
		addr := Config.WorkerAddrs[i]
		if httputil.IsWorkerAvailableWithSwitch(addr, Config.CoordinatorAddrs[0]) {
			workers = append(workers, addr)
			if len(workers) == lim {
				return workers, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("Not enough workers, asked %d got %d", lim, len(workers)))
}
