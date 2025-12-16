package server

import (
	"log"
	"time"

	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
	"github.com/Croohand/mapreduce/common/timeutil"
)

func monitorCoordinators() {
	for {
		timeutil.Sleep(time.Second * 3)
		anyActive := false
		for _, addr := range Config.CoordinatorAddrs[1:] {
			var status responses.CoordinatorStatus
			resp, err := httpClient.Get(addr + "/IsAlive")
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
			log.Println(Config.Name + " switch state to active, become active coordinator")
			state = "active"
			Config.LastJournalTs = time.Now()
			go RunServices()
			return
		}
	}
}
