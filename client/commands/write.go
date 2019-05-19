package commands

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
)

func tryValidateBlocks(txId string, blocks []fsutil.BlockInfoEx) error {
	for l := 0; l < len(blocks); l += 50 {
		r := l + 50
		if r > len(blocks) {
			r = len(blocks)
		}
		b, err := json.Marshal(blocks[l:r])
		if err != nil {
			return err
		}
		resp, err := http.PostForm(mrConfig.Host+"/Transaction/ValidateBlocks", url.Values{"TransactionId": {txId}, "Blocks": {string(b)}})
		if err != nil {
			return err
		}
		if err := httputil.GetError(resp); err != nil {
			return err
		}
	}
	return nil
}

func Write(path string, doAppend bool) {
	if !fsutil.ValidateFilePath(path) {
		log.Panic("Invalid file path " + path)
	}
	txId, txHandler := startWriteTransaction([]string{path})
	defer txHandler.close()

	blocks := []fsutil.BlockInfoEx{}

	in := make(chan string, 10)
	done := make(chan error)
	go httputil.WriteBlocks(in, mrConfig.Host, txId, mrConfig.MrConfig, &blocks, done)
	for scanner := bufio.NewScanner(os.Stdin); scanner.Scan(); {
		line := scanner.Text()
		select {
		case err := <-done:
			log.Panic(err)
		case in <- line:
		}
	}
	close(in)
	err := <-done
	if err != nil {
		log.Panic(err)
	}

	var blockIds fsutil.BlockIds
	for _, block := range blocks {
		blockIds = append(blockIds, block.Id)
	}

	for l := 0; l < len(blockIds); l += 500 {
		r := l + 500
		if r > len(blockIds) {
			r = len(blockIds)
		}
		if l > 0 {
			doAppend = true
		}
		resp, err := http.PostForm(mrConfig.Host+"/File/Write", url.Values{"Path": {path}, "Append": {strconv.FormatBool(doAppend)}, "BlockIds": blockIds[l:r]})
		if err != nil {
			log.Panic("Failed to write path in database: " + err.Error())
		}
		if err := httputil.GetError(resp); err != nil {
			log.Panic("Failed to write path in database: " + err.Error())
		}
	}

	err = tryValidateBlocks(txId, blocks)

	all := map[string]struct{}{}
	for _, block := range blocks {
		for _, slave := range block.Slaves {
			all[slave] = struct{}{}
		}
	}
	for slave := range all {
		if !httputil.IsSlaveAvailable(slave) {
			continue
		}
		resp, err := http.PostForm(slave+"/Transaction/Remove", url.Values{"TransactionId": {txId}})
		if err != nil {
			continue
		}
		httputil.GetError(resp)
	}

	if err != nil {
		log.Panic("Failed to close write transaction: " + err.Error())
	}
}
