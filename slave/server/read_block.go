package server

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Croohand/mapreduce/common/fsutil"
)

func readBlock(id string, w io.Writer) error {
	if !checkBlock(id, "").Exists {
		return errors.New(fmt.Sprintf("Block with id %s doesn't exist", id))
	}
	file, err := os.Open(fsutil.GetBlockPath(id, ""))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(w, file)
	return err
}
