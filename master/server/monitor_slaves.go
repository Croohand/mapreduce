package server

import (
	"log"
	"time"

	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/timeutil"
	"github.com/Croohand/mapreduce/common/wrrors"
	"github.com/Croohand/mapreduce/master/server/dbase"
)

var (
	lastEnsured   = map[string]time.Time{}
	lastAvailable = map[string]time.Time{}
)

func ensureSlave(slave string) error {
	blocks, err := dbase.GetKeys(slave, "")
	if err != nil {
		return err
	}
	for _, blockId := range blocks {
		go ensureBlock(blockId)
	}
	return nil
}

func monitorSlaves() {
	wrr := wrrors.New("monitorSlaves")
	for {
		timeutil.Sleep(time.Minute)
		log.Println(Config.Name + " monitorSlaves starting new iteration")
		for _, slave := range Config.SlaveAddrs {
			lastAv, hasAv := lastAvailable[slave]
			lastEn, hasEn := lastEnsured[slave]
			if !httputil.IsSlaveAvailable(slave) {
				if (!hasAv || time.Since(lastAv) > 5*time.Minute) && (!hasEn || time.Since(lastEn) > time.Since(lastAv)) {
					log.Println(wrr.WrapS("Slave " + slave + " seems to be down"))
					err := ensureSlave(slave)
					if err != nil {
						log.Println(wrr.Wrap(err))
						continue
					}
					lastEnsured[slave] = time.Now()
				}
				continue
			}
			lastAvailable[slave] = time.Now()
			if hasEn && (!hasAv || time.Since(lastAv) > time.Since(lastEn)) {
				log.Println(wrr.WrapS("Slave " + slave + " is up again"))
				err := ensureSlave(slave)
				if err != nil {
					log.Println(wrr.Wrap(err))
					continue
				}
			}
		}
	}
}
