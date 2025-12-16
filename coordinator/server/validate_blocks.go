package server

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
	"github.com/Croohand/mapreduce/coordinator/server/dbase"
)

func validateBlocks(txId string, blocks []*fsutil.BlockInfoEx) error {
	for _, block := range blocks {
		available := false
		for _, worker := range block.Workers {
			if !httputil.IsWorkerAvailable(worker) {
				continue
			}
			resp, err := httpClient.PostForm(worker+"/Block/Check", url.Values{"BlockId": {block.Id}, "TransactionId": {txId}})
			if err != nil {
				log.Println(err)
				continue
			}
			var bStatus responses.BlockStatus
			if err := httputil.GetJson(resp, &bStatus); err != nil {
				log.Println(err)
				continue
			}
			if bStatus.Exists {
				available = true
			}
		}
		if !available {
			return errors.New(fmt.Sprintf("Can't validate blocks, block with id %s is absent", block.Id))
		}
	}
	for _, block := range blocks {
		written := []string{}
		for _, worker := range block.Workers {
			if !httputil.IsWorkerAvailable(worker) {
				continue
			}
			resp, err := httpClient.PostForm(worker+"/Block/Validate", url.Values{"BlockId": {block.Id}, "TransactionId": {txId}})
			if err != nil {
				log.Println(err)
				continue
			}
			if err := httputil.GetError(resp); err != nil {
				log.Println(err)
				continue
			}
			written = append(written, worker)
			err = dbase.Set(worker, block.Id, []byte("1"))
			if err != nil {
				log.Println(err)
				continue
			}
		}
		block.Workers = written
		err := dbase.SetObject(dbase.Blocks, block.Id, block.BlockInfo)
		if err != nil {
			log.Println(err)
			continue
		}
		go ensureBlock(block.Id)
	}
	return nil
}
