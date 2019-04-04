package server

import (
	"encoding/json"
	"net/http"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/wrrors"
)

func getArgs(w http.ResponseWriter, r *http.Request, wrr wrrors.Wrror) (id, path string, got bool) {
	id, path = r.PostFormValue("TransactionId"), r.PostFormValue("Path")
	got = false
	if !fsutil.ValidateTransactionId(id) {
		http.Error(w, wrr.SWrap("invalid transaction id "+id), http.StatusBadRequest)
		return
	}
	if !fsutil.ValidateFilePath(path) {
		http.Error(w, wrr.SWrap("invalid file path "+path), http.StatusBadRequest)
		return
	}
	got = true
	return
}

func updateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("updateTransactionHandler")
	id, path, got := getArgs(w, r, wrr)
	if !got {
		return
	}
	resp, err := updateTransaction(id, path)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}

func isAliveTransactionHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("isAliveTransactionHandler")
	id, path, got := getArgs(w, r, wrr)
	if !got {
		return
	}
	httputil.WriteResponse(w, isAliveTransaction(id, path), nil)
}

func validateWriteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("validateWriteTransactionHandler")
	id, path, got := getArgs(w, r, wrr)
	if !got {
		return
	}
	blocks := r.PostFormValue("PathInfo")
	var pathInfo fsutil.PathInfo
	if err := json.Unmarshal([]byte(blocks), &pathInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := validateWriteTransaction(id, path, pathInfo, []byte(blocks))
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}
