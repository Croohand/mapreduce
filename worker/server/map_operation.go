package server

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/Croohand/mapreduce/common/fsutil"
)

func mapOperation(blockId, txId string, reducers int) error {
	mainPath := filepath.Join(fsutil.GetTxDir(txId), "bin", "main")
	blockPath := fsutil.GetBlockPath(blockId, "")
	resPath := fsutil.GetBlockPath(blockId, txId) + "-res"
	err := os.Mkdir(resPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	mapCmd := exec.Command(mainPath, "Map", blockPath, resPath, strconv.Itoa(reducers))
	var stderr bytes.Buffer
	mapCmd.Stderr = &stderr
	err = mapCmd.Run()
	if err != nil {
		return errors.New(err.Error() + ": " + stderr.String())
	}
	return nil
}
