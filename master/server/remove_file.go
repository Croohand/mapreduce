package server

import (
	"log"

	"github.com/Croohand/mapreduce/master/server/dbase"
)

func removeFile(path string) error {
	err := dbase.Del(dbase.PathTxs, path)
	if err != nil {
		log.Println(err)
	}
	return dbase.Del(dbase.Files, path)
}
