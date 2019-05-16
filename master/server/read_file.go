package server

import (
	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/responses"
	"github.com/Croohand/mapreduce/master/server/dbase"
)

func readFile(path string) (blocks responses.PathBlocks, err error) {
	var blockIds fsutil.BlockIds
	err = dbase.GetObject(dbase.Files, path, &blockIds)
	if err != nil {
		return
	}
	for _, blockId := range blockIds {
		var blockInfo fsutil.BlockInfo
		err = dbase.GetObject(dbase.Blocks, blockId, &blockInfo)
		if err != nil {
			return
		}
		blocks = append(blocks, fsutil.BlockInfoEx{blockId, blockInfo})
	}
	return
}
