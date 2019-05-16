package server

import (
	"log"
	"net/http"
	"net/url"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/wrrors"
	"github.com/Croohand/mapreduce/master/server/dbase"
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
	newSlaves := []string{}
	for _, slave := range block.Slaves {
		newSlaves = append(newSlaves, slave)
		if !httputil.IsSlaveAvailable(slave) {
			continue
		}
		resp, err := http.PostForm(slave+"/Block/Remove", url.Values{"BlockId": {blockId}})
		if err != nil {
			log.Println(wrr.Wrap(err))
			continue
		}
		if err := httputil.GetError(resp); err != nil {
			log.Println(wrr.Wrap(err))
			continue
		}
		if err := dbase.Del(slave, blockId); err != nil {
			log.Println(wrr.Wrap(err))
		}
		newSlaves = newSlaves[:len(newSlaves)-1]
	}
	block.Slaves = newSlaves
	if len(block.Slaves) > 0 {
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
