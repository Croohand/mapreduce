package server

import (
	"io"
	"mime/multipart"
	"os"

	"github.com/Croohand/mapreduce/common/fsutil"
)

func writeBlock(id, txId string, file multipart.File) error {
	err := fsutil.CreateTxDir(txId)
	if err != nil {
		return err
	}
	path := fsutil.GetBlockPath(id, txId)
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	return err
}
