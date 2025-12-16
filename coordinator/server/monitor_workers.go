package server

import (
	"log"
	"time"

	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/timeutil"
	"github.com/Croohand/mapreduce/common/wrrors"
	"github.com/Croohand/mapreduce/coordinator/server/dbase"
)

var (
	lastEnsured   = map[string]time.Time{}
	lastAvailable = map[string]time.Time{}
)

func ensureWorker(worker string) error {
	blocks, err := dbase.GetKeys(worker, "")
	if err != nil {
		return err
	}
	for _, blockId := range blocks {
		go ensureBlock(blockId)
	}
	return nil
}

func monitorWorkers() {
	wrr := wrrors.New("monitorWorkers")
	for {
		log.Println(Config.Name + " monitorWorkers starting new iteration")
		for _, worker := range Config.WorkerAddrs {
			lastAv, hasAv := lastAvailable[worker]
			lastEn, hasEn := lastEnsured[worker]
			if !httputil.IsWorkerAvailableWithSwitch(worker, Config.CoordinatorAddrs[0]) {
				if (!hasAv || time.Since(lastAv) > 5*time.Minute) && (!hasEn || time.Since(lastEn) > time.Since(lastAv)) {
					log.Println(wrr.WrapS("Worker " + worker + " seems to be down"))
					err := ensureWorker(worker)
					if err != nil {
						log.Println(wrr.Wrap(err))
						continue
					}
					lastEnsured[worker] = time.Now()
				}
				continue
			}
			lastAvailable[worker] = time.Now()
			if hasEn && (!hasAv || time.Since(lastAv) > time.Since(lastEn)) {
				log.Println(wrr.WrapS("Worker " + worker + " is up again"))
				err := ensureWorker(worker)
				if err != nil {
					log.Println(wrr.Wrap(err))
					continue
				}
			}
		}

		for _, scheduler := range Config.SchedulerAddrs {
			httputil.IsWorkerAvailableWithSwitch(scheduler, Config.CoordinatorAddrs[0])
		}

		timeutil.Sleep(time.Minute)
	}
}
