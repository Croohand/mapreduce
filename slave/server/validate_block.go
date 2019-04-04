package server

import (
	"os"
	"path/filepath"

	"github.com/Croohand/mapreduce/common/responses"
)

func validateBlock(id, transaction string) (*responses.Answer, error) {
	oldPath := filepath.Join("transactions", transaction, id)
	newPath := filepath.Join("files", id)
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return nil, err
	}
	return &responses.Answer{true}, nil
}
