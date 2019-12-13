package server

import (
	"log"
	"net/url"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/wrrors"
	"github.com/Croohand/mapreduce/master/server/dbase"
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
	for _, slave := range block.Slaves {
		if httputil.IsSlaveAvailable(slave) {
			available[slave] = struct{}{}
		} else {
			unavailable = append(unavailable, slave)
		}
	}
	if getMrConfig().ReplicationFactor > len(available) {
		for _, slave := range Config.SlaveAddrs {
			_, has := available[slave]
			if !has && httputil.IsSlaveAvailable(slave) {
				for from := range available {
					resp, err := httpClient.PostForm(from+"/Block/Copy", url.Values{"BlockId": {blockId}, "Where": {slave}})
					if err != nil {
						log.Println(wrr.Wrap(err))
						continue
					}
					if err := httputil.GetError(resp); err != nil {
						log.Println(wrr.Wrap(err))
						continue
					}
					available[slave] = struct{}{}
					err = dbase.Set(slave, blockId, []byte("1"))
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
		for slave := range available {
			avCopy[slave] = struct{}{}
		}
		for slave := range avCopy {
			if len(available) <= getMrConfig().ReplicationFactor {
				break
			}
			resp, err := httpClient.PostForm(slave+"/Block/Remove", url.Values{"BlockId": {blockId}})
			if err != nil {
				log.Println(wrr.Wrap(err))
				continue
			}
			if err := httputil.GetError(resp); err != nil {
				log.Println(wrr.Wrap(err))
				continue
			}
			delete(available, slave)
			err = dbase.Del(slave, blockId)
			if err != nil {
				log.Println(wrr.Wrap(err))
			}
		}
	}
	for _, slave := range unavailable {
		available[slave] = struct{}{}
	}
	newSlaves := []string{}
	for slave := range available {
		newSlaves = append(newSlaves, slave)
	}
	block.Slaves = newSlaves
	return wrr.Wrap(dbase.SetObject(dbase.Blocks, blockId, block))
}
