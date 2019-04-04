package server

import (
	"os"
	"path/filepath"

	"github.com/Croohand/mapreduce/common/responses"
)

func checkBlock(id, transaction string) *responses.BlockStatus {
	path := filepath.Join("files", id)
	if transaction != "" {
		path = filepath.Join("transactions", transaction, id)
	}
	_, err := os.Stat(path)
	return &responses.BlockStatus{!os.IsNotExist(err)}
}
