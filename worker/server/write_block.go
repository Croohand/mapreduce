package server

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/Croohand/mapreduce/common/fsutil"
)

func writeBlock(id, txId string, file multipart.File, shuffle bool) error {
	err := fsutil.CreateTxDir(txId)
	if err != nil {
		return err
	}
	path := fsutil.GetBlockPath(id, txId)
	if shuffle {
		err := fsutil.CreateShuffleDir(txId)
		if err != nil {
			return err
		}
		path = filepath.Join(fsutil.GetTxDir(txId), "shuffle", id)
	}
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	return err
}
