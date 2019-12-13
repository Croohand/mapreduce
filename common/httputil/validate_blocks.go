package httputil

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"

	"github.com/Croohand/mapreduce/common/fsutil"
)

func TryWritePath(mrHost, txId, path string, blocks []fsutil.BlockInfoEx, doAppend bool) error {
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
		resp, err := httpClient.PostForm(mrHost+"/File/Write", url.Values{"Path": {path}, "Append": {strconv.FormatBool(doAppend)}, "BlockIds": blockIds[l:r]})
		if err != nil {
			return errors.New("Failed to write path in database: " + err.Error())
		}
		if err := GetError(resp); err != nil {
			return errors.New("Failed to write path in database: " + err.Error())
		}
	}

	return nil
}

func TryValidateBlocks(mrHost, txId string, blocks []fsutil.BlockInfoEx) error {
	for l := 0; l < len(blocks); l += 50 {
		r := l + 50
		if r > len(blocks) {
			r = len(blocks)
		}
		b, err := json.Marshal(blocks[l:r])
		if err != nil {
			return err
		}
		resp, err := httpClient.PostForm(mrHost+"/Transaction/ValidateBlocks", url.Values{"TransactionId": {txId}, "Blocks": {string(b)}})
		if err != nil {
			return err
		}
		if err := GetError(resp); err != nil {
			return err
		}
	}
	return nil
}

func CleanUp(txId string, blocks []fsutil.BlockInfoEx) {
	all := map[string]struct{}{}
	for _, block := range blocks {
		for _, slave := range block.Slaves {
			all[slave] = struct{}{}
		}
	}
	for slave := range all {
		if !IsSlaveAvailable(slave) {
			continue
		}
		resp, err := httpClient.PostForm(slave+"/Transaction/Remove", url.Values{"TransactionId": {txId}})
		if err != nil {
			continue
		}
		GetError(resp)
	}
}
