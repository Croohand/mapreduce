package server

import (
	"net/http"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/wrrors"
)

func removeTransactionHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("removeTransactionHandler")
	id := r.PostFormValue("TransactionId")
	if !fsutil.ValidateTransactionId(id) {
		http.Error(w, wrr.SWrap("invalid transaction id "+id), http.StatusBadRequest)
		return
	}
	resp, err := removeTransaction(id)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}
