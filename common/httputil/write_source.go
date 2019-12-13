package httputil

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func WriteSource(dir, name, txId, where string) error {
	path := filepath.Join(dir, name)
	var b bytes.Buffer

	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormField("SourceName")
	if err != nil {
		return err
	}
	if _, err = fw.Write([]byte(name)); err != nil {
		return err
	}

	fw, err = w.CreateFormField("TransactionId")
	if err != nil {
		return err
	}
	if _, err = fw.Write([]byte(txId)); err != nil {
		return err
	}

	fw, err = w.CreateFormFile("SourceCode", "File")
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
	w.Close()

	req, err := http.NewRequest("POST", where+"/Source/Write", bytes.NewBuffer(b.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Close = true
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	return GetError(resp)
}

func WriteSources(srcsPath, txId, where string) error {
	files, err := ioutil.ReadDir(srcsPath)
	if err != nil {
		return err
	}
	any := false
	for _, fileInfo := range files {
		if !fileInfo.IsDir() {
			err := WriteSource(srcsPath, fileInfo.Name(), txId, where)
			if err != nil {
				return err
			}
			any = true
		}
	}

	if !any {
		return errors.New("No source files sent")
	}
	return nil
}
