package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/wrrors"
)

func mapOperationHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("mapOperationHandler")
	blockId := r.PostFormValue("BlockId")
	if !fsutil.ValidateBlockId(blockId) {
		http.Error(w, wrr.SWrap("Invalid block id "+blockId), http.StatusBadRequest)
		return
	}

	txId := r.PostFormValue("TransactionId")
	if !fsutil.ValidateTransactionId(txId) || len(txId) == 0 {
		http.Error(w, wrr.SWrap("Invalid transaction id "+txId), http.StatusBadRequest)
		return
	}

	reducers, err := strconv.Atoi(r.PostFormValue("Reducers"))
	if err != nil {
		http.Error(w, wrr.SWrap(err.Error()), http.StatusBadRequest)
		return
	}
	if reducers <= 0 {
		http.Error(w, wrr.SWrap("Invalid number of reducers "+strconv.Itoa(reducers)), http.StatusBadRequest)
		return
	}

	err = mapOperation(blockId, txId, reducers)
	httputil.WriteResponse(w, nil, wrr.Wrap(err))
}

func getStatusOperationHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("getStatusOperationHandler")
	txId := r.PostFormValue("TransactionId")
	if !fsutil.ValidateTransactionId(txId) || len(txId) == 0 {
		http.Error(w, wrr.SWrap("Invalid transaction id "+txId), http.StatusBadRequest)
		return
	}

	resp, err := getStatusOperation(txId)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}

func reduceOperationHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("reduceOperationHandler")
	txId := r.PostFormValue("TransactionId")
	if !fsutil.ValidateTransactionId(txId) || len(txId) == 0 {
		http.Error(w, wrr.SWrap("Invalid transaction id "+txId), http.StatusBadRequest)
		return
	}

	resp, err := reduceOperation(txId)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}

func sendResultsOperationHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("sendResultsOperationHandler")
	txId := r.PostFormValue("TransactionId")
	if !fsutil.ValidateTransactionId(txId) || len(txId) == 0 {
		http.Error(w, wrr.SWrap("Invalid transaction id "+txId), http.StatusBadRequest)
		return
	}

	blockId := r.PostFormValue("BlockId")
	if !fsutil.ValidateBlockId(blockId) {
		http.Error(w, wrr.SWrap("Invalid block id "+blockId), http.StatusBadRequest)
		return
	}

	where := r.PostForm["Where"]
	dst := map[string]string{}
	for _, raw := range where {
		toks := strings.SplitN(raw, " ", 2)
		if len(toks) != 2 {
			http.Error(w, wrr.SWrap("Invalid destinations field"), http.StatusBadRequest)
			return
		}
		dst[toks[0]] = toks[1]
	}

	err := sendResultsOperation(blockId, txId, dst)
	httputil.WriteResponse(w, nil, wrr.Wrap(err))
}

func prepareMapReduceOperationHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("prepareMapReduceOperationHandler")
	out := r.PostFormValue("Out")
	if !fsutil.ValidateFilePath(out) {
		http.Error(w, wrr.SWrap("Invalid file path "+out), http.StatusBadRequest)
		return
	}

	in := r.PostForm["In"]
	for _, path := range in {
		if !fsutil.ValidateFilePath(path) {
			http.Error(w, wrr.SWrap("Invalid file path "+path), http.StatusBadRequest)
			return
		}
	}

	resp, err := prepareMapReduceOperation(in, out)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}

func startMapReduceOperationHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("startMapReduceOperationHandler")
	out := r.PostFormValue("Out")
	if !fsutil.ValidateFilePath(out) {
		http.Error(w, wrr.SWrap("Invalid file path "+out), http.StatusBadRequest)
		return
	}

	in := r.PostForm["In"]
	for _, path := range in {
		if !fsutil.ValidateFilePath(path) {
			http.Error(w, wrr.SWrap("Invalid file path "+path), http.StatusBadRequest)
			return
		}
	}

	readTxId := r.PostFormValue("ReadTransactionId")
	if !fsutil.ValidateTransactionId(readTxId) || len(readTxId) == 0 {
		http.Error(w, wrr.SWrap("Invalid transaction id "+readTxId), http.StatusBadRequest)
		return
	}

	writeTxId := r.PostFormValue("WriteTransactionId")
	if !fsutil.ValidateTransactionId(writeTxId) || len(writeTxId) == 0 {
		http.Error(w, wrr.SWrap("Invalid transaction id "+writeTxId), http.StatusBadRequest)
		return
	}

	reducers, err := strconv.Atoi(r.PostFormValue("Reducers"))
	if err != nil {
		http.Error(w, wrr.SWrap(err.Error()), http.StatusBadRequest)
		return
	}
	if reducers <= 0 {
		http.Error(w, wrr.SWrap("Invalid number of reducers "+strconv.Itoa(reducers)), http.StatusBadRequest)
		return
	}

	mappers, err := strconv.Atoi(r.PostFormValue("Mappers"))
	if err != nil {
		http.Error(w, wrr.SWrap(err.Error()), http.StatusBadRequest)
		return
	}
	if mappers <= 0 {
		http.Error(w, wrr.SWrap("Invalid number of mappers "+strconv.Itoa(mappers)), http.StatusBadRequest)
		return
	}

	go startMapReduceOperation(in, out, readTxId, writeTxId, mappers, reducers)
	httputil.WriteResponse(w, nil, nil)
}
