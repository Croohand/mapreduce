package server

import (
	"log"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/timeutil"
	"github.com/Croohand/mapreduce/common/wrrors"
	"github.com/Croohand/mapreduce/master/server/dbase"
)

func monitorFiles() {
	wrr := wrrors.New("monitorFiles")
	for {
		timeutil.HugeSleep()
		log.Println(Config.Name + " monitorFiles starting new iteration")
		files, err := dbase.GetKeys(dbase.Files, "")
		if err != nil {
			log.Println(wrr.Wrap(err))
			continue
		}
		used := map[string]struct{}{}
		fail := false
		for _, file := range files {
			var blockIds fsutil.BlockIds
			err := dbase.GetObject(dbase.Files, file, &blockIds)
			if err != nil {
				log.Println(wrr.Wrap(err))
				fail = true
				break
			}
			for _, blockId := range blockIds {
				used[blockId] = struct{}{}
			}
		}
		if !fail {
			blockIds, err := dbase.GetKeys(dbase.Blocks, "")
			if err != nil {
				log.Println(wrr.Wrap(err))
				continue
			}
			for _, blockId := range blockIds {
				_, need := used[blockId]
				if need {
					err := ensureBlock(blockId)
					if err != nil {
						log.Println(wrr.Wrap(err))
						continue
					}
				} else {
					err := removeBlock(blockId)
					if err != nil {
						log.Println(wrr.Wrap(err))
						continue
					}
				}
			}
		}
	}
}
