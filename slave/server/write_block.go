package server

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/Croohand/mapreduce/common/responses"
)

func writeBlock(id, transaction, meta string, file multipart.File) (*responses.Answer, error) {
	path := filepath.Join("transactions", transaction, id)
	finalPath := filepath.Join("files", id)
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return nil, errors.New("block with path " + path + " already exists")
	}
	_, err = os.Stat(finalPath)
	if !os.IsNotExist(err) {
		return nil, errors.New("block with final path " + finalPath + " already exists")
	}
	txPath := filepath.Join("transactions", transaction)
	_, err = os.Stat(txPath)
	if os.IsNotExist(err) {
		if err := os.Mkdir(txPath, os.ModePerm); err != nil {
			return nil, err
		}
		metaPath := filepath.Join("transactions", transaction, "meta")
		dst, err := os.OpenFile(metaPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return nil, err
		}
		defer dst.Close()
		dst.Write([]byte(meta))
	}
	dst, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, err
	}
	return &responses.Answer{true}, nil
}
