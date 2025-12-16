package server

import (
	"net/http"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/wrrors"
)

func removeTransactionHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("removeTransactionHandler")
	txId := r.PostFormValue("TransactionId")
	if !fsutil.ValidateTransactionId(txId) {
		http.Error(w, wrr.SWrap("Invalid transaction id "+txId), http.StatusBadRequest)
		return
	}
	err := removeTransaction(txId)
	httputil.WriteResponse(w, nil, wrr.Wrap(err))
}
