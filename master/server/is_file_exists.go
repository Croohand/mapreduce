package server

import (
	"github.com/Croohand/mapreduce/common/responses"
	"github.com/Croohand/mapreduce/master/server/dbase"
)

func isFileExists(path string) (*responses.FileStatus, error) {
	has, err := dbase.Has(dbase.Files, path)
	return &responses.FileStatus{has}, err
}
