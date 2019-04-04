package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
)

func sendBlock(txId string, cur bytes.Buffer, path string, lower, upper int) (block *fsutil.BlockInfo, err error) {
	resp, err := http.Get(mrConfig.Host + "/GetAvailableSlaves")
	if err != nil {
		return
	}
	var slaves []string
	if err = httputil.GetJson(resp, &slaves); err != nil {
		return
	}
	block = &fsutil.BlockInfo{Id: fsutil.GenerateBlockId(), Lower: lower, Upper: upper, Slaves: slaves}
	var b bytes.Buffer

	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormField("BlockId")
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

	fw, err = w.CreateFormField("Meta")
	if err != nil {
		return
	}
	tmp, err := json.Marshal(struct{ Path string }{path})
	if err != nil {
		return
	}
	if _, err = fw.Write(tmp); err != nil {
		return
	}

	fw, err = w.CreateFormFile("File", "Block")
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, &cur); err != nil {
		return
	}

	w.Close()

	for _, addr := range block.Slaves {
		req, err := http.NewRequest("POST", addr+"/Block/Write", bytes.NewBuffer(b.Bytes()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		var ans responses.Answer
		if err := httputil.GetJson(resp, &ans); err != nil {
			return nil, err
		}
	}
	return
}

func Write(path string) {
	if !fsutil.ValidateFilePath(path) {
		log.Fatal("invalid file path " + path)
	}
	txId := startTransaction(path)
	var blocks fsutil.PathInfo
	var cur bytes.Buffer

	lower, upper := 0, 0

	trySendBlock := func() {
		block, err := sendBlock(txId, cur, path, lower, upper)
		if err != nil {
			log.Fatal("send block failed with error %s", err.Error())
		}
		blocks = append(blocks, *block)
		lower = upper
		cur.Reset()
	}

	for scanner := bufio.NewScanner(os.Stdin); scanner.Scan(); {
		line := scanner.Text()
		if _, err := cur.WriteString(line); err != nil {
			log.Fatal(err)
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
	b, err := json.Marshal(blocks)
	if err != nil {
		log.Fatal("failed to validate write transaction " + err.Error())
	}
	resp, err := http.PostForm(mrConfig.Host+"/Transaction/ValidateWrite", url.Values{"Path": {path}, "TransactionId": {txId}, "PathInfo": {string(b)}})
	if err != nil {
		log.Fatal("failed to validate write transaction " + err.Error())
	}
	var ans responses.Answer
	if err := httputil.GetJson(resp, &ans); err != nil {
		log.Fatal("failed to validate write transaction " + err.Error())
	}
}
