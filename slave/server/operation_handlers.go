package server

import (
	"net/http"
	"strconv"

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
