package server

import (
	"os"
	"path/filepath"

	"github.com/Croohand/mapreduce/common/responses"
)

func removeTransaction(id string) (*responses.Answer, error) {
	path := filepath.Join("transactions", id)
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		if err := os.RemoveAll(path); err != nil {
			return nil, err
		}
	}
	return &responses.Answer{true}, nil
}
