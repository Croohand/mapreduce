package server

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/Croohand/mapreduce/common/httputil"
)

func getAvailableSlaves(lim int) ([]string, error) {
	slaves := make([]string, 0)
	for _, i := range rand.Perm(len(Config.SlaveAddrs)) {
		addr := Config.SlaveAddrs[i]
		if httputil.IsSlaveAvailableWithSwitch(addr, Config.MasterAddrs[0]) {
			slaves = append(slaves, addr)
			if len(slaves) == lim {
				return slaves, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("Not enough slaves, asked %d got %d", lim, len(slaves)))
}
