package server

import (
	"log"

	"github.com/Croohand/mapreduce/common/timeutil"
	"github.com/Croohand/mapreduce/common/wrrors"
	"github.com/Croohand/mapreduce/master/server/dbase"
)

func monitorBlocks() {
	wrr := wrrors.New("monitorBlocks")
	for {
		timeutil.HugeSleep()
		log.Println(Config.Name + " monitorBlocks starting new iteration")
		blockIds, err := dbase.GetKeys(dbase.Blocks, "")
		if err != nil {
			log.Println(wrr.Wrap(err))
			continue
		}
		for _, blockId := range blockIds {
			err := ensureBlock(blockId)
			if err != nil {
				log.Println(wrr.Wrap(err))
				continue
			}
		}
	}
}
