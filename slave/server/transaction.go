package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
)

func removeTransactionInner(id string) error {
	path := filepath.Join("transactions", id)
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

func removeTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("TransactionId")
	if !fsutil.ValidateTransactionId(id) {
		http.Error(w, "invalid transaction id "+id, http.StatusBadRequest)
		return
	}
	if err := removeTransactionInner(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httputil.WriteJson(w, struct{ Success bool }{true})
}
