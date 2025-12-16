package server

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/Croohand/mapreduce/common/fsutil"
)

func sendResultsOperation(blockId, txId string, dst map[string]string) error {
	resultsDir := fsutil.GetBlockPath(blockId, txId) + "-res"
	files, err := ioutil.ReadDir(resultsDir)
	if err != nil {
		return err
	}
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			return errors.New("Directory found in " + resultsDir)
		}
		where, has := dst[fileInfo.Name()]
		if !has {
			continue
		}
		var err error
		for try := 0; try < 3; try++ {
			err = copyFile(filepath.Join(resultsDir, fileInfo.Name()), blockId, txId, where, true)
			if err == nil {
				break
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
