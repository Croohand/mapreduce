package httputil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/responses"
)

func sendBlock(mrHost, txId string, cur bytes.Buffer, lower, upper int) (block *fsutil.BlockInfoEx, err error) {
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

	resp, err := http.Get(mrHost + "/GetAvailableSlaves")
	if err != nil {
		return
	}
	var slaves []string
	if err = GetJson(resp, &slaves); err != nil {
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
	if err := GetError(resp); err != nil {
		return nil, err
	}
	for _, addr := range block.Slaves[1:] {
		resp, err = http.PostForm(initialSlave+"/Block/Copy", url.Values{"TransactionId": {txId}, "BlockId": {block.Id}, "Where": {addr}})
		if err != nil {
			return nil, err
		}
		if err := GetError(resp); err != nil {
			return nil, err
		}
	}
	return
}

func WriteBlocks(in <-chan string, mrHost, txId string, config responses.MrConfig, blocks *[]fsutil.BlockInfoEx, done chan<- error) {
	var cur bytes.Buffer

	lower, upper := 0, 0

	trySendBlock := func() error {
		var err error
		var block *fsutil.BlockInfoEx
		for try := 0; try < 3; try++ {
			block, err = sendBlock(mrHost, txId, cur, lower, upper)
			if err == nil {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		if err != nil {
			return errors.New(fmt.Sprintf("Send block failed with error %s", err.Error()))
		}
		*blocks = append(*blocks, *block)
		lower = upper
		cur.Reset()
		return nil
	}

	for line := range in {
		if len(line) == 0 {
			continue
		}
		if len(line) >= config.MaxRowLength {
			done <- errors.New(fmt.Sprintf("Max row length %d exceeded", config.MaxRowLength))
			return
		}
		if _, err := cur.WriteString(line); err != nil {
			done <- err
			return
		}
		upper += 1
		cur.WriteByte('\n')
		if cur.Len() >= config.BlockSize {
			err := trySendBlock()
			if err != nil {
				done <- err
				return
			}
		}
	}

	if cur.Len() > 0 {
		err := trySendBlock()
		done <- err
	} else {
		done <- nil
	}
}
