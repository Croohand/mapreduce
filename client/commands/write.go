package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
)

func sendBlock(txId string, cur bytes.Buffer, path string, lower, upper int) (block *fsutil.BlockInfoEx, err error) {
	block = &fsutil.BlockInfoEx{fsutil.GenerateBlockId(), fsutil.BlockInfo{Lower: lower, Upper: upper}}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("Block", "File")
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, &cur); err != nil {
		return
	}

	fw, err = w.CreateFormField("BlockId")
	if err != nil {
		return
	}
	if _, err = fw.Write([]byte(block.Id)); err != nil {
		return
	}

	fw, err = w.CreateFormField("TransactionId")
	if err != nil {
		return
	}
	if _, err = fw.Write([]byte(txId)); err != nil {
		return
	}

	w.Close()

	resp, err := http.Get(mrConfig.Host + "/GetAvailableSlaves")
	if err != nil {
		return
	}
	var slaves []string
	if err = httputil.GetJson(resp, &slaves); err != nil {
		return
	}
	block.Slaves = slaves
	initialSlave := block.Slaves[0]
	req, err := http.NewRequest("POST", initialSlave+"/Block/Write", &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Close = true
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := httputil.GetError(resp); err != nil {
		return nil, err
	}
	for _, addr := range block.Slaves[1:] {
		resp, err = http.PostForm(initialSlave+"/Block/Copy", url.Values{"TransactionId": {txId}, "BlockId": {block.Id}, "Where": {addr}})
		if err != nil {
			return nil, err
		}
		if err := httputil.GetError(resp); err != nil {
			return nil, err
		}
	}
	return
}

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
	var blocks []fsutil.BlockInfoEx
	var cur bytes.Buffer

	lower, upper := 0, 0

	trySendBlock := func() {
		var err error
		var block *fsutil.BlockInfoEx
		for try := 0; try < 3; try++ {
			block, err = sendBlock(txId, cur, path, lower, upper)
			if err == nil {
				break
			}
			time.Sleep(time.Second)
		}
		if err != nil {
			log.Panic(fmt.Sprintf("Send block failed with error %s", err.Error()))
		}
		blocks = append(blocks, *block)
		lower = upper
		cur.Reset()
	}

	for scanner := bufio.NewScanner(os.Stdin); scanner.Scan(); {
		line := scanner.Text()
		if len(line) >= mrConfig.MaxRowLength {
			log.Panic(fmt.Sprintf("Max row length %d exceeded", mrConfig.MaxRowLength))
		}
		if _, err := cur.WriteString(line); err != nil {
			log.Panic(err)
		}
		upper += 1
		cur.WriteByte('\n')
		if cur.Len() >= mrConfig.BlockSize {
			trySendBlock()
		}
	}

	if cur.Len() > 0 {
		trySendBlock()
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

	err := tryValidateBlocks(txId, blocks)

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
