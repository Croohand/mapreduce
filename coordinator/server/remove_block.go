package server

import (
	"log"
	"net/url"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/wrrors"
	"github.com/Croohand/mapreduce/coordinator/server/dbase"
)

func removeBlock(blockId string) error {
	wrr := wrrors.New("removeBlock")
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
	newWorkers := []string{}
	for _, worker := range block.Workers {
		newWorkers = append(newWorkers, worker)
		if !httputil.IsWorkerAvailable(worker) {
			continue
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
		if err := dbase.Del(worker, blockId); err != nil {
			log.Println(wrr.Wrap(err))
		}
		newWorkers = newWorkers[:len(newWorkers)-1]
	}
	block.Workers = newWorkers
	if len(block.Workers) > 0 {
		err := dbase.SetObject(dbase.Blocks, blockId, block)
		if err != nil {
			return wrr.Wrap(err)
		}
	} else {
		err := dbase.Del(dbase.Blocks, blockId)
		if err != nil {
			return wrr.Wrap(err)
		}
	}
	return nil
}
