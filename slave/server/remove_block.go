package server

import (
	"os"
	"path/filepath"

	"github.com/Croohand/mapreduce/common/responses"
)

func removeBlock(id string) (*responses.Answer, error) {
	path := filepath.Join("files", id)
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		if err := os.Remove(path); err != nil {
			return nil, err
		}
	}
	return &responses.Answer{true}, nil
}
