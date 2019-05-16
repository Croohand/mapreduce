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

func copyBlock(id, txId, where string) error {
	if !checkBlock(id, txId).Exists {
		return errors.New(fmt.Sprintf("Block %s doesn't exist", fsutil.GetBlockPath(id, txId)))
	}
	var b bytes.Buffer

	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormField("BlockId")
	if err != nil {
		return err
	}
	if _, err = fw.Write([]byte(id)); err != nil {
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

	path := fsutil.GetBlockPath(id, txId)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = io.Copy(fw, file); err != nil {
		return err
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
	if err := httputil.GetError(resp); err != nil {
		return err
	}

	return nil
}
