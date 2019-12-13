package httputil

import (
	"net/url"

	"github.com/Croohand/mapreduce/common/responses"
)

func IsSlaveAvailable(addr string) bool {
	resp, err := avChecker.Get(addr + "/IsAlive")
	if err == nil {
		var status responses.SlaveStatus
		if GetJson(resp, &status) == nil && status.Alive {
			return true
		}
	}
	return false
}

func IsSlaveAvailableWithSwitch(addr, master string) bool {
	resp, err := avChecker.PostForm(addr+"/IsAlive", url.Values{"Switch": {"true"}, "MasterAddr": {master}})
	if err == nil {
		var status responses.SlaveStatus
		if GetJson(resp, &status) == nil && status.Alive {
			return true
		}
	}
	return false
}
