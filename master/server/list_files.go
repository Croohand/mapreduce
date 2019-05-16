package server

import (
	"github.com/Croohand/mapreduce/common/responses"
	"github.com/Croohand/mapreduce/master/server/dbase"
)

func listFiles(prefix string) (responses.ListedFiles, error) {
	return dbase.GetKeys(dbase.Files, prefix)
}
