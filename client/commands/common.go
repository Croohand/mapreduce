package commands

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
)

func startTransaction(path string) string {
	started := make(chan bool, 1)
	id := fsutil.GenerateBlockId()
	go startTransactionInner(path, id, started)
	<-started
	return id
}

func startTransactionInner(path string, id string, started chan bool) {
	const pingInterval = 10 * time.Second
	firstPing := false
	for {
		resp, err := http.PostForm(mrConfig.Host+"/Transaction/Update", url.Values{"TransactionId": []string{id}, "Path": []string{path}})
		if err != nil {
			log.Fatal(err)
		}
		var success struct{ Success bool }
		if err := httputil.GetJson(resp, &success); err != nil {
			log.Fatal(err)
		}
		if !success.Success {
			log.Fatal("couldn't take transaction")
		}
		if !firstPing {
			firstPing = true
			started <- true
		}
		time.Sleep(pingInterval)
	}
}
