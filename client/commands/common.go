package commands

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
)

func startTransaction(path string) string {
	started := make(chan bool, 1)
	id := fsutil.GenerateTransactionId()
	go startTransactionInner(path, id, started)
	<-started
	return id
}

func startTransactionInner(path string, id string, started chan bool) {
	const pingInterval = 10 * time.Second
	firstPing := false
	for {
		resp, err := http.PostForm(mrConfig.Host+"/Transaction/Update", url.Values{"TransactionId": {id}, "Path": {path}})
		if err != nil {
			log.Fatal(err)
		}
		var ans responses.Answer
		if err := httputil.GetJson(resp, &ans); err != nil {
			log.Fatal(err)
		}
		if !ans.Success {
			log.Fatal("couldn't take transaction")
		}
		if !firstPing {
			firstPing = true
			started <- true
		}
		time.Sleep(pingInterval)
	}
}
