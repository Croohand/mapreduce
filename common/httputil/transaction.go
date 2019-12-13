package httputil

import (
	"log"
	"net/url"
	"time"

	"github.com/Croohand/mapreduce/common/responses"
)

type TransactionHandler struct {
	signalExit chan bool
	signalDone chan bool
}

func (txHandler *TransactionHandler) Close() {
	txHandler.signalExit <- true
	<-txHandler.signalDone
}

func NewTxHandler() *TransactionHandler {
	return &TransactionHandler{make(chan bool), make(chan bool)}
}

func startTransactionInner(mrHost string, paths []string, txType string, ping bool) (string, *TransactionHandler) {
	resp, err := httpClient.PostForm(mrHost+"/Transaction/Start", url.Values{"Paths": paths, "TransactionType": {txType}})
	if err != nil {
		log.Panic(err)
	}
	var txInfo responses.StartedTransaction
	if err := GetJson(resp, &txInfo); err != nil {
		log.Panic(err)
	}
	if !ping {
		return txInfo.Id, nil
	}
	txHandler := NewTxHandler()
	go PingTransaction(mrHost, txInfo.Id, txHandler)
	return txInfo.Id, txHandler
}

func StartReadTransaction(mrHost string, paths []string, ping bool) (string, *TransactionHandler) {
	return startTransactionInner(mrHost, paths, "Read", ping)
}

func StartWriteTransaction(mrHost string, paths []string, ping bool) (string, *TransactionHandler) {
	return startTransactionInner(mrHost, paths, "Write", ping)
}

func PingTransaction(mrHost, txId string, txHandler *TransactionHandler) {
	const pingInterval = 10 * time.Second
	ticker := time.NewTicker(pingInterval)
	for {
		select {
		case <-txHandler.signalExit:
			resp, err := httpClient.PostForm(mrHost+"/Transaction/Close", url.Values{"TransactionId": {txId}})
			if err != nil {
				log.Println("Failed to close transaction: " + err.Error())
			}
			if err := GetError(resp); err != nil {
				log.Println("Failed to close transaction: " + err.Error())
			}
			txHandler.signalDone <- true
			return
		case <-ticker.C:
			resp, err := httpClient.PostForm(mrHost+"/Transaction/Update", url.Values{"TransactionId": {txId}})
			if err != nil {
				log.Panic(err)
			}
			if err := GetError(resp); err != nil {
				log.Panic(err)
			}
		}
	}
}
