package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
	"github.com/Croohand/mapreduce/master/server/dbase"
)

func validateBlocks(txId string, blocks []*fsutil.BlockInfoEx) error {
	for _, block := range blocks {
		available := false
		for _, slave := range block.Slaves {
			if !httputil.IsSlaveAvailable(slave) {
				continue
			}
			resp, err := http.PostForm(slave+"/Block/Check", url.Values{"BlockId": {block.Id}, "TransactionId": {txId}})
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
		for _, slave := range block.Slaves {
			if !httputil.IsSlaveAvailable(slave) {
				continue
			}
			resp, err := http.PostForm(slave+"/Block/Validate", url.Values{"BlockId": {block.Id}, "TransactionId": {txId}})
			if err != nil {
				log.Println(err)
				continue
			}
			if err := httputil.GetError(resp); err != nil {
				log.Println(err)
				continue
			}
			written = append(written, slave)
			err = dbase.Set(slave, block.Id, []byte("1"))
			if err != nil {
				log.Println(err)
				continue
			}
		}
		block.Slaves = written
		err := dbase.SetObject(dbase.Blocks, block.Id, block.BlockInfo)
		if err != nil {
			log.Println(err)
			continue
		}
		go ensureBlock(block.Id)
	}
	return nil
}
