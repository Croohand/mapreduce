package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
)

func copyFile(path, blockId, txId, where string, shuffle bool) error {
	var b bytes.Buffer

	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormField("BlockId")
	if err != nil {
		return err
	}
	if _, err = fw.Write([]byte(blockId)); err != nil {
		return err
	}

	fw, err = w.CreateFormField("TransactionId")
	if err != nil {
		return err
	}
	if _, err = fw.Write([]byte(txId)); err != nil {
		return err
	}

	fw, err = w.CreateFormFile("Block", "File")
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = io.Copy(fw, file); err != nil {
		return err
	}

	if shuffle {
		fw, err := w.CreateFormField("Shuffle")
		if err != nil {
			return err
		}
		_, err = fw.Write([]byte("true"))
		if err != nil {
			return err
		}
	}
	w.Close()

	req, err := http.NewRequest("POST", where+"/Block/Write", bytes.NewBuffer(b.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Close = true
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return httputil.GetError(resp)
}

func copyBlock(id, txId, where string) error {
	path := fsutil.GetBlockPath(id, txId)
	if !checkBlock(id, txId).Exists {
		return errors.New(fmt.Sprintf("Block %s doesn't exist", path))
	}
	return copyFile(path, id, txId, where, false)
}
