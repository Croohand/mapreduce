package server

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Croohand/mapreduce/common/blockutil"
	"github.com/Croohand/mapreduce/common/httputil"
)

func checkBlock(w http.ResponseWriter, r *http.Request) {
	id, transaction := r.PostFormValue("BlockId"), r.PostFormValue("TransactionId")
	if !blockutil.ValidateId(id) {
		http.Error(w, "invalid block id "+id, http.StatusBadRequest)
		return
	}
	if transaction != "" && !blockutil.ValidateId(transaction) {
		http.Error(w, "invalid transaction id "+transaction, http.StatusBadRequest)
		return
	}
	path := filepath.Join("files", id)
	if transaction != "" {
		path = filepath.Join("transactions", transaction, id)
	}
	_, err := os.Stat(path)
	httputil.WriteJson(w, struct{ Exists bool }{!os.IsNotExist(err)})
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	id, transaction, meta := r.PostFormValue("BlockId"), r.PostFormValue("TransactionId"), r.PostFormValue("Meta")
	if !blockutil.ValidateId(id) {
		http.Error(w, "invalid block id "+id, http.StatusBadRequest)
		return
	}
	if !blockutil.ValidateId(transaction) {
		http.Error(w, "invalid transaction id "+transaction, http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("File")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	txPath := filepath.Join("transactions", transaction)
	_, err = os.Stat(txPath)
	if os.IsNotExist(err) {
		if err := os.Mkdir(txPath, os.ModePerm); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		metaPath := filepath.Join("transactions", transaction, "meta")
		dst, err := os.OpenFile(metaPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		dst.Write([]byte(meta))
	}
	path := filepath.Join("transactions", transaction, id)
	finalPath := filepath.Join("files", id)
	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		http.Error(w, "block with path "+path+" already exists", http.StatusBadRequest)
		return
	}
	_, err = os.Stat(finalPath)
	if !os.IsNotExist(err) {
		http.Error(w, "block with final path "+finalPath+" already exists", http.StatusBadRequest)
		return
	}
	dst, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httputil.WriteJson(w, struct{ Success bool }{true})
}

func validateBlock(w http.ResponseWriter, r *http.Request) {
	id, transaction := r.PostFormValue("BlockId"), r.PostFormValue("TransactionId")
	if !blockutil.ValidateId(id) {
		http.Error(w, "invalid block id "+id, http.StatusBadRequest)
		return
	}
	if !blockutil.ValidateId(transaction) {
		http.Error(w, "invalid transaction id "+id, http.StatusBadRequest)
		return
	}
	oldPath := filepath.Join("transactions", transaction, id)
	newPath := filepath.Join("files", id)
	err := os.Rename(oldPath, newPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	httputil.WriteJson(w, struct{ Success bool }{true})
}

func removeBlock(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("BlockId")
	if !blockutil.ValidateId(id) {
		http.Error(w, "invalid block id "+id, http.StatusBadRequest)
		return
	}
	path := filepath.Join("files", id)
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		if err := os.Remove(path); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	httputil.WriteJson(w, struct{ Success bool }{true})
}

func readBlock(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("BlockId")
	if !blockutil.ValidateId(id) {
		http.Error(w, "invalid block id "+id, http.StatusBadRequest)
		return
	}
	path := filepath.Join("files", id)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		http.Error(w, "block with id "+id+" doesn't exist", http.StatusBadRequest)
		return
	}
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
