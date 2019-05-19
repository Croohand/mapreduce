package server

import (
	"bytes"
	"errors"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
)

func reduceOperation(txId string) (responses.PathBlocks, error) {
	var mrConfig responses.MrConfig
	resp, err := http.Get(Config.MasterAddr + "/GetMrConfig")
	if err != nil {
		return nil, errors.New("Couldn't get MR config from master, error: " + err.Error())
	}
	if err := httputil.GetJson(resp, &mrConfig); err != nil {
		return nil, errors.New("Couldn't get MR config from master, error: " + err.Error())
	}
	mainPath := filepath.Join(fsutil.GetTxDir(txId), "bin", "main")
	dirPath := filepath.Join(fsutil.GetTxDir(txId), "shuffle")
	reduceCmd := exec.Command(mainPath, "Reduce", dirPath)
	var stderr bytes.Buffer
	reduceCmd.Stderr = &stderr
	p, err := reduceCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = reduceCmd.Start()
	if err != nil {
		return nil, err
	}
	blocks, err := httputil.WriteBlocks(p, Config.MasterAddr, txId, mrConfig)
	if err == nil {
		err = reduceCmd.Wait()
		if err != nil {
			err = errors.New(err.Error() + ": " + stderr.String())
		}
	} else {
		reduceCmd.Wait()
	}
	if err != nil {
		return nil, err
	}
	return blocks, err
}
