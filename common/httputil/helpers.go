package httputil

import (
	"net/url"

	"github.com/Croohand/mapreduce/common/responses"
)

func IsWorkerAvailable(addr string) bool {
	resp, err := avChecker.Get(addr + "/IsAlive")
	if err == nil {
		var status responses.WorkerStatus
		if GetJson(resp, &status) == nil && status.Alive {
			return true
		}
	}
	return false
}

func IsWorkerAvailableWithSwitch(addr, coordinator string) bool {
	resp, err := avChecker.PostForm(addr+"/IsAlive", url.Values{"Switch": {"true"}, "CoordinatorAddr": {coordinator}})
	if err == nil {
		var status responses.WorkerStatus
		if GetJson(resp, &status) == nil && status.Alive {
			return true
		}
	}
	return false
}
