package server

import (
	"io"
	"mime/multipart"
	"os"

	"github.com/Croohand/mapreduce/common/fsutil"
)

func writeSource(name, txId string, file multipart.File) error {
	err := fsutil.CreateSourcesDir(txId)
	if err != nil {
		return err
	}
	path := fsutil.GetSourcePath(name, txId)
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	return err
}
