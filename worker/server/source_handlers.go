package server

import (
	"net/http"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/wrrors"
)

func writeSourceHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("writeSourceHandler")
	name, txId := r.PostFormValue("SourceName"), r.PostFormValue("TransactionId")
	if !fsutil.ValidateTransactionId(txId) {
		http.Error(w, wrr.SWrap("Invalid transaction id "+txId), http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("SourceCode")
	if err != nil {
		http.Error(w, wrr.Wrap(err).Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	err = writeSource(name, txId, file)
	httputil.WriteResponse(w, nil, wrr.Wrap(err))
}

func buildSourceHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("buildSourceHandler")
	txId := r.PostFormValue("TransactionId")
	if !fsutil.ValidateTransactionId(txId) {
		http.Error(w, wrr.SWrap("Invalid transaction id "+txId), http.StatusBadRequest)
		return
	}
	err := buildSource(txId)
	httputil.WriteResponse(w, nil, wrr.Wrap(err))
}
