package server

import (
	"os"

	"github.com/Croohand/mapreduce/common/fsutil"
)

func removeBlock(id string) error {
	path := fsutil.GetBlockPath(id, "")
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}
