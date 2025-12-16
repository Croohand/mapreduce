package server

import (
	"os"

	"github.com/Croohand/mapreduce/common/fsutil"
)

func validateBlock(id, txId string) error {
	oldPath := fsutil.GetBlockPath(id, txId)
	newPath := fsutil.GetBlockPath(id, "")
	return os.Rename(oldPath, newPath)
}
