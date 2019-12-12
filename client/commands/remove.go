package commands

import (
	"log"
	"net/http"
	"net/url"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
)

func Remove(path string) {
	if !existsInner(path) {
		log.Panic("File path " + path + " doesn't exist")
	}
	if !fsutil.ValidateFilePath(path) {
		log.Panic("Invalid file path " + path)
	}

	_, txHandler := startWriteTransaction([]string{path})
	defer txHandler.Close()

	resp, err := http.PostForm(mrConfig.GetHost()+"/File/Remove", url.Values{"Path": {path}})
	if err != nil {
		log.Panic(err)
	}
	if err := httputil.GetError(resp); err != nil {
		log.Panic(err)
	}
}
