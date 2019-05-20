package commands

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
)

func MapReduce(in []string, out, srcsPath string, mappers, reducers int, detached bool) {
	for _, path := range in {
		if !existsInner(path) {
			log.Panic("File path " + path + " doesn't exist")
		}
	}
	if !fsutil.ValidateFilePath(out) {
		log.Panic("Invalid file path " + out)
	}

	if reducers <= 0 {
		log.Panic("Invalid number of reducers")
	}

	if mappers <= 0 {
		log.Panic("Invalid number of mappers")
	}

	if !strings.HasSuffix(srcsPath, "mruserlib") {
		log.Panic("Path to user library needs to have suffix mruserlib")
	}

	stat, err := os.Stat(srcsPath)
	if err != nil && os.IsNotExist(err) {
		log.Panic(err)
	}

	if !stat.IsDir() {
		log.Panic("Path to user library is not a directory")
	}

	resp, err := http.PostForm(mrConfig.Host+"/GetAvailableScheduler", url.Values{})
	if err != nil {
		log.Panic(err)
	}
	scheduler := ""
	if err := httputil.GetJson(resp, &scheduler); err != nil {
		log.Panic(err)
	}

	var txs responses.PreparedOperation
	resp, err = http.PostForm(scheduler+"/Operation/PrepareMapReduce", url.Values{"In": in, "Out": {out}})
	if err != nil {
		log.Panic(err)
	}
	if err := httputil.GetJson(resp, &txs); err != nil {
		log.Panic(err)
	}
	txId := txs.WriteTxId

	err = httputil.WriteSources(srcsPath, txId, scheduler)
	if err != nil {
		log.Panic(err)
	}

	resp, err = http.PostForm(scheduler+"/Operation/StartMapReduce", url.Values{"In": in, "Out": {out}, "ReadTransactionId": {txs.ReadTxId}, "WriteTransactionId": {txs.WriteTxId}, "Reducers": {strconv.Itoa(reducers)}, "Mappers": {strconv.Itoa(mappers)}})
	if err != nil {
		log.Panic(err)
	}
	if err := httputil.GetJson(resp, &txs); err != nil {
		log.Panic(err)
	}
	if detached {
		return
	}
}
