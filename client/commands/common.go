package commands

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Croohand/mapreduce/common/blockutil"
	"github.com/Croohand/mapreduce/common/httputil"
)

func startTransaction(path string) string {
	started := make(chan bool, 1)
	id := blockutil.GenerateId()
	go startTransactionInner(path, id, started)
	<-started
	return id
}

func startTransactionInner(path string, id string, started chan bool) {
	const pingInterval = 10 * time.Second
	const pingRetries = 10
	firstPing := false
	for {
		for i := 0; i < pingRetries; i++ {
			resp, err := http.PostForm(mrConfig.Host+"/Transaction/Update", url.Values{"TransactionId": []string{id}, "Path": []string{path}})
			if err == nil {
				var success struct{ Success bool }
				if err := httputil.GetJson(resp, &success); err == nil {
					if !success.Success {
						log.Fatal("couldn't take transaction")
					}
					if !firstPing {
						firstPing = true
						started <- true
					}
					break
				} else {
					log.Println(err)
				}
			} else {
				log.Println(err)
			}
			if i == pingRetries-1 {
				log.Fatal("couldn't take transaction, max retries exceeded")
			}
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(pingInterval)
	}
}
