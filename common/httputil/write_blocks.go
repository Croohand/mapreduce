package httputil

import (
	"bufio"
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
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Close = true
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	if err = GetError(resp); err != nil {
		return
	}
	for _, addr := range block.Slaves[1:] {
		resp, err = http.PostForm(initialSlave+"/Block/Copy", url.Values{"TransactionId": {txId}, "BlockId": {block.Id}, "Where": {addr}})
		if err != nil {
			return
		}
		if err = GetError(resp); err != nil {
			return
		}
	}
	return
}

func WriteBlocks(in io.Reader, mrHost, txId string, config responses.MrConfig) ([]fsutil.BlockInfoEx, error) {
	blocks := []fsutil.BlockInfoEx{}
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
		blocks = append(blocks, *block)
		lower = upper
		cur.Reset()
		return nil
	}

	sc := bufio.NewScanner(in)
	for sc.Scan() {
		line := sc.Text()
		if len(line) == 0 {
			continue
		}
		if len(line) >= config.MaxRowLength {
			return nil, errors.New(fmt.Sprintf("Max row length %d exceeded", config.MaxRowLength))
		}
		if _, err := cur.WriteString(line); err != nil {
			return nil, err
		}
		upper += 1
		cur.WriteByte('\n')
		if cur.Len() >= config.BlockSize {
			err := trySendBlock()
			if err != nil {
				return nil, err
			}
		}
	}
	if sc.Err() != nil {
		return nil, sc.Err()
	}

	if cur.Len() > 0 {
		err := trySendBlock()
		return blocks, err
	} else {
		return blocks, nil
	}
}
