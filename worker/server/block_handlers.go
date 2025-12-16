package server

import (
	"net/http"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/wrrors"
)

func checkBlockHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("checkBlockHandler")
	id, txId := r.PostFormValue("BlockId"), r.PostFormValue("TransactionId")
	if !fsutil.ValidateBlockId(id) {
		http.Error(w, wrr.SWrap("Invalid block id "+id), http.StatusBadRequest)
		return
	}
	if !fsutil.ValidateTransactionId(txId) {
		http.Error(w, wrr.SWrap("Invalid transaction id "+txId), http.StatusBadRequest)
		return
	}
	httputil.WriteResponse(w, checkBlock(id, txId), nil)
}

func writeBlockHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("writeBlockHandler")
	id, txId, shuffle := r.PostFormValue("BlockId"), r.PostFormValue("TransactionId"), r.PostFormValue("Shuffle")
	if !fsutil.ValidateBlockId(id) {
		http.Error(w, wrr.SWrap("Invalid block id "+id), http.StatusBadRequest)
		return
	}
	if !fsutil.ValidateTransactionId(txId) {
		http.Error(w, wrr.SWrap("Invalid transaction id "+txId), http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("Block")
	if err != nil {
		http.Error(w, wrr.Wrap(err).Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	err = writeBlock(id, txId, file, (shuffle == "true"))
	httputil.WriteResponse(w, nil, wrr.Wrap(err))
}

func copyBlockHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("copyBlockHandler")
	id, txId, where := r.PostFormValue("BlockId"), r.PostFormValue("TransactionId"), r.PostFormValue("Where")
	if !fsutil.ValidateBlockId(id) {
		http.Error(w, wrr.SWrap("Invalid block id "+id), http.StatusBadRequest)
		return
	}
	if !fsutil.ValidateTransactionId(txId) {
		http.Error(w, wrr.SWrap("Invalid transaction id "+txId), http.StatusBadRequest)
		return
	}
	err := copyBlock(id, txId, where)
	httputil.WriteResponse(w, nil, wrr.Wrap(err))
}

func validateBlockHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("validateBlockHandler")
	id, txId := r.PostFormValue("BlockId"), r.PostFormValue("TransactionId")
	if !fsutil.ValidateBlockId(id) {
		http.Error(w, wrr.SWrap("Invalid block id "+id), http.StatusBadRequest)
		return
	}
	if !fsutil.ValidateTransactionId(txId) {
		http.Error(w, wrr.SWrap("Invalid transaction id "+id), http.StatusBadRequest)
		return
	}
	err := validateBlock(id, txId)
	httputil.WriteResponse(w, nil, wrr.Wrap(err))
}

func removeBlockHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("removeBlockHandler")
	id := r.PostFormValue("BlockId")
	if !fsutil.ValidateBlockId(id) {
		http.Error(w, wrr.SWrap("Invalid block id "+id), http.StatusBadRequest)
		return
	}
	err := removeBlock(id)
	httputil.WriteResponse(w, nil, wrr.Wrap(err))
}

func readBlockHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("readBlockHandler")
	id := r.PostFormValue("BlockId")
	if !fsutil.ValidateBlockId(id) {
		http.Error(w, wrr.SWrap("Invalid block id "+id), http.StatusBadRequest)
		return
	}
	err := readBlock(id, w)
	if err != nil {
		httputil.WriteResponse(w, nil, wrr.Wrap(err))
	}
}
