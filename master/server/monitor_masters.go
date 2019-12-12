package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
	"github.com/Croohand/mapreduce/common/timeutil"
)

func monitorMasters() {
	timeutil.Sleep(time.Second * 5)
	for {
		anyActive := false
		for _, addr := range Config.MasterAddrs[1:] {
			var status responses.MasterStatus
			resp, err := http.Get(addr + "/IsAlive")
			if err != nil {
				continue
			}
			if err = httputil.GetJson(resp, &status); err != nil {
				continue
			}
			if status.State == "active" {
				anyActive = true
			}
		}
		if !anyActive {
			log.Println(Config.Name + " switch state to active, become active master")
			state = "active"
			go RunServices()
			return
		}
		timeutil.Sleep(time.Second * 3)
	}
}
