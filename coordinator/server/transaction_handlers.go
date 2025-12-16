package server

import (
	"encoding/json"
	"net/http"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/wrrors"
)

func getCommonArgs(w http.ResponseWriter, r *http.Request, wrr wrrors.Wrror) (txId string, got bool) {
	txId = r.PostFormValue("TransactionId")
	got = false
	if !fsutil.ValidateTransactionId(txId) {
		http.Error(w, wrr.SWrap("Invalid transaction id "+txId), http.StatusBadRequest)
		return
	}
	got = true
	return
}

func updateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("updateTransactionHandler")
	txId, got := getCommonArgs(w, r, wrr)
	if !got {
		return
	}
	err := updateTransaction(txId)
	httputil.WriteResponse(w, nil, wrr.Wrap(err))
}

func startTransactionHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("startTransactionHandler")
	txTypeRaw := r.PostFormValue("TransactionType")
	var txType fsutil.TransactionType
	switch txTypeRaw {
	case "Read":
		txType = fsutil.TxTypeRead
	case "Write":
		txType = fsutil.TxTypeWrite
	default:
		http.Error(w, wrr.SWrap("Invalid transaction type"), http.StatusBadRequest)
		return
	}
	paths := r.PostForm["Paths"]
	for _, path := range paths {
		if !fsutil.ValidateFilePath(path) {
			http.Error(w, wrr.SWrap("Invalid file path "+path), http.StatusBadRequest)
			return
		}
	}
	resp, err := startTransaction(paths, txType)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}

func closeTransactionHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("closeTransactionHandler")
	txId, got := getCommonArgs(w, r, wrr)
	if !got {
		return
	}
	err := closeTransaction(txId)
	httputil.WriteResponse(w, nil, wrr.Wrap(err))
}

func isAliveTransactionHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("isAliveTransactionHandler")
	txId, got := getCommonArgs(w, r, wrr)
	if !got {
		return
	}
	resp, err := isAliveTransaction(txId)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}

func validateBlocksHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("validateBlocksHandler")
	txId, got := getCommonArgs(w, r, wrr)
	if !got {
		return
	}
	blocksRaw := r.PostFormValue("Blocks")
	var blocks []*fsutil.BlockInfoEx
	if err := json.Unmarshal([]byte(blocksRaw), &blocks); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := validateBlocks(txId, blocks)
	httputil.WriteResponse(w, nil, wrr.Wrap(err))
}
