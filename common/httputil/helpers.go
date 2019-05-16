package httputil

import (
	"net/http"
	"time"

	"github.com/Croohand/mapreduce/common/responses"
)

var avChecker = http.Client{Timeout: time.Duration(100 * time.Millisecond)}

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
