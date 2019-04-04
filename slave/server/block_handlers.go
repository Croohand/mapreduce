package server

import (
	"net/http"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/wrrors"
)

func checkBlockHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("checkBlockHandler")
	id, transaction := r.PostFormValue("BlockId"), r.PostFormValue("TransactionId")
	if !fsutil.ValidateBlockId(id) {
		http.Error(w, wrr.SWrap("invalid block id "+id), http.StatusBadRequest)
		return
	}
	if transaction != "" && !fsutil.ValidateTransactionId(transaction) {
		http.Error(w, wrr.SWrap("invalid transaction id "+transaction), http.StatusBadRequest)
		return
	}
	httputil.WriteResponse(w, checkBlock(id, transaction), nil)
}

func writeBlockHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("writeBlockHandler")
	id, transaction, meta := r.PostFormValue("BlockId"), r.PostFormValue("TransactionId"), r.PostFormValue("Meta")
	if !fsutil.ValidateBlockId(id) {
		http.Error(w, wrr.SWrap("invalid block id "+id), http.StatusBadRequest)
		return
	}
	if !fsutil.ValidateTransactionId(transaction) {
		http.Error(w, wrr.SWrap("invalid transaction id "+transaction), http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("File")
	if err != nil {
		http.Error(w, wrr.Wrap(err).Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	resp, err := writeBlock(id, transaction, meta, file)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}

func validateBlockHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("validateBlockHandler")
	id, transaction := r.PostFormValue("BlockId"), r.PostFormValue("TransactionId")
	if !fsutil.ValidateBlockId(id) {
		http.Error(w, wrr.SWrap("invalid block id "+id), http.StatusBadRequest)
		return
	}
	if !fsutil.ValidateTransactionId(transaction) {
		http.Error(w, wrr.SWrap("invalid transaction id "+id), http.StatusBadRequest)
		return
	}
	resp, err := validateBlock(id, transaction)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}

func removeBlockHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("removeBlockHandler")
	id := r.PostFormValue("BlockId")
	if !fsutil.ValidateBlockId(id) {
		http.Error(w, wrr.SWrap("invalid block id "+id), http.StatusBadRequest)
		return
	}
	resp, err := removeBlock(id)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}

func readBlockHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("readBlockHandler")
	id := r.PostFormValue("BlockId")
	if !fsutil.ValidateBlockId(id) {
		http.Error(w, wrr.SWrap("invalid block id "+id), http.StatusBadRequest)
		return
	}
	err := readBlock(id, w)
	if err != nil {
		httputil.WriteResponse(w, nil, wrr.Wrap(err))
	}
}
