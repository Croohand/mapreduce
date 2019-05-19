package commands

import (
	"log"
	"net/http"
	"net/url"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
)

func Merge(in []string, out string) {
	for _, path := range in {
		if !existsInner(path) {
			log.Panic("File path " + path + " doesn't exist")
		}
	}
	if !fsutil.ValidateFilePath(out) {
		log.Panic("Invalid file path " + out)
	}

	_, readTxHandler := startReadTransaction(in)
	defer readTxHandler.Close()
	_, writeTxHandler := startWriteTransaction([]string{out})
	defer writeTxHandler.Close()

	resp, err := http.PostForm(mrConfig.Host+"/File/Merge", url.Values{"In": in, "Out": {out}})
	if err != nil {
		log.Panic(err)
	}
	if err := httputil.GetError(resp); err != nil {
		log.Panic(err)
	}
}
