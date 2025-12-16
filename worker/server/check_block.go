package server

import (
	"os"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/responses"
)

func checkBlock(id, txId string) *responses.BlockStatus {
	path := fsutil.GetBlockPath(id, txId)
	_, err := os.Stat(path)
	return &responses.BlockStatus{!os.IsNotExist(err)}
}
