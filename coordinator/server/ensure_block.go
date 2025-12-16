package server

import (
	"log"
	"net/url"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/wrrors"
	"github.com/Croohand/mapreduce/coordinator/server/dbase"
)

func ensureBlock(blockId string) error {
	wrr := wrrors.New("ensureBlock")
	has, err := dbase.Has(dbase.Blocks, blockId)
	if err != nil {
		return wrr.Wrap(err)
	}
	if !has {
		return wrr.WrapS("No block with id " + blockId)
	}
	var block fsutil.BlockInfo
	err = dbase.GetObject(dbase.Blocks, blockId, &block)
	if err != nil {
		return wrr.Wrap(err)
	}
	available := map[string]struct{}{}
	unavailable := []string{}
	for _, worker := range block.Workers {
		if httputil.IsWorkerAvailable(worker) {
			available[worker] = struct{}{}
		} else {
			unavailable = append(unavailable, worker)
		}
	}
	if getMrConfig().ReplicationFactor > len(available) {
		for _, worker := range Config.WorkerAddrs {
			_, has := available[worker]
			if !has && httputil.IsWorkerAvailable(worker) {
				for from := range available {
					resp, err := httpClient.PostForm(from+"/Block/Copy", url.Values{"BlockId": {blockId}, "Where": {worker}})
					if err != nil {
						log.Println(wrr.Wrap(err))
						continue
					}
					if err := httputil.GetError(resp); err != nil {
						log.Println(wrr.Wrap(err))
						continue
					}
					available[worker] = struct{}{}
					err = dbase.Set(worker, blockId, []byte("1"))
					if err != nil {
						log.Println(wrr.Wrap(err))
					}
					break
				}
				if len(available) >= getMrConfig().ReplicationFactor {
					break
				}
			}
		}
	} else {
		avCopy := make(map[string]struct{})
		for worker := range available {
			avCopy[worker] = struct{}{}
		}
		for worker := range avCopy {
			if len(available) <= getMrConfig().ReplicationFactor {
				break
			}
			resp, err := httpClient.PostForm(worker+"/Block/Remove", url.Values{"BlockId": {blockId}})
			if err != nil {
				log.Println(wrr.Wrap(err))
				continue
			}
			if err := httputil.GetError(resp); err != nil {
				log.Println(wrr.Wrap(err))
				continue
			}
			delete(available, worker)
			err = dbase.Del(worker, blockId)
			if err != nil {
				log.Println(wrr.Wrap(err))
			}
		}
	}
	for _, worker := range unavailable {
		available[worker] = struct{}{}
	}
	newWorkers := []string{}
	for worker := range available {
		newWorkers = append(newWorkers, worker)
	}
	block.Workers = newWorkers
	return wrr.Wrap(dbase.SetObject(dbase.Blocks, blockId, block))
}
