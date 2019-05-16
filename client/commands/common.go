package commands

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
)

type transactionHandler struct {
	signalExit chan bool
	signalDone chan bool
}

func (txHandler *transactionHandler) close() {
	txHandler.signalExit <- true
	<-txHandler.signalDone
}

func startTransactionInner(paths []string, txType string) (string, transactionHandler) {
	resp, err := http.PostForm(mrConfig.Host+"/Transaction/Start", url.Values{"Paths": paths, "TransactionType": {txType}})
	if err != nil {
		log.Panic(err)
	}
	var txInfo responses.StartedTransaction
	if err := httputil.GetJson(resp, &txInfo); err != nil {
		log.Panic(err)
	}
	txHandler := transactionHandler{make(chan bool), make(chan bool)}
	go pingTransaction(txInfo.Id, txHandler)
	return txInfo.Id, txHandler
}

func startReadTransaction(paths []string) (string, transactionHandler) {
	return startTransactionInner(paths, "Read")
}

func startWriteTransaction(paths []string) (string, transactionHandler) {
	return startTransactionInner(paths, "Write")
}

func pingTransaction(txId string, txHandler transactionHandler) {
	const pingInterval = 10 * time.Second
	ticker := time.NewTicker(pingInterval)
	for {
		select {
		case <-txHandler.signalExit:
			resp, err := http.PostForm(mrConfig.Host+"/Transaction/Close", url.Values{"TransactionId": {txId}})
			if err != nil {
				log.Println("Failed to close transaction: " + err.Error())
			}
			if err := httputil.GetError(resp); err != nil {
				log.Println("Failed to close transaction: " + err.Error())
			}
			txHandler.signalDone <- true
			return
		case <-ticker.C:
			resp, err := http.PostForm(mrConfig.Host+"/Transaction/Update", url.Values{"TransactionId": {txId}})
			if err != nil {
				log.Panic(err)
			}
			if err := httputil.GetError(resp); err != nil {
				log.Panic(err)
			}
		}
	}
}
