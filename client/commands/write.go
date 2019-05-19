package commands

import (
	"log"
	"os"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
)

func Write(path string, doAppend bool) {
	if !fsutil.ValidateFilePath(path) {
		log.Panic("Invalid file path " + path)
	}
	txId, txHandler := startWriteTransaction([]string{path})
	defer txHandler.Close()

	blocks, err := httputil.WriteBlocks(os.Stdin, mrConfig.Host, txId, mrConfig.MrConfig)

	if err != nil {
		log.Panic(err)
	}

	err = httputil.TryWritePath(mrConfig.Host, txId, path, blocks, doAppend)

	if err != nil {
		log.Panic(err)
	}

	err = httputil.TryValidateBlocks(mrConfig.Host, txId, blocks)
	httputil.CleanUp(txId, blocks)

	if err != nil {
		log.Panic("Failed to close write transaction: " + err.Error())
	}
}
