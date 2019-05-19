package server

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/Croohand/mapreduce/common/httputil"
)

func getAvailableReducers(lim int) ([]string, error) {
	reducers := make([]string, 0)
	for _, i := range rand.Perm(len(Config.SlaveAddrs)) {
		addr := Config.SlaveAddrs[i]
		if httputil.IsSlaveAvailable(addr) {
			reducers = append(reducers, addr)
			if len(reducers) == lim {
				return reducers, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("Not enough reducers, asked %d got %d", lim, len(reducers)))
}
