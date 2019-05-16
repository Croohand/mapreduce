package server

import (
	"net/http"
	"strconv"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/wrrors"
)

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("listFilesHandler")
	prefix := r.PostFormValue("Prefix")
	resp, err := listFiles(prefix)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}

func isFileExistsHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("isFileExistsHandler")
	path := r.PostFormValue("Path")
	if !fsutil.ValidateFilePath(path) {
		http.Error(w, wrr.SWrap("Invalid file path "+path), http.StatusBadRequest)
		return
	}
	resp, err := isFileExists(path)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}

func removeFileHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("removeFileHandler")
	path := r.PostFormValue("Path")
	if !fsutil.ValidateFilePath(path) {
		http.Error(w, wrr.SWrap("Invalid file path "+path), http.StatusBadRequest)
		return
	}
	err := removeFile(path)
	httputil.WriteResponse(w, nil, wrr.Wrap(err))
}

func readFileHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("readFileHandler")
	path := r.PostFormValue("Path")
	if !fsutil.ValidateFilePath(path) {
		http.Error(w, wrr.SWrap("Invalid file path "+path), http.StatusBadRequest)
		return
	}
	resp, err := readFile(path)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}

func writeFileHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("writeFileHandler")
	path := r.PostFormValue("Path")
	if !fsutil.ValidateFilePath(path) {
		http.Error(w, wrr.SWrap("Invalid file path "+path), http.StatusBadRequest)
		return
	}
	app, err := strconv.ParseBool(r.PostFormValue("Append"))
	if err != nil {
		http.Error(w, wrr.SWrap("Couldn't parse bool with error: "+err.Error()), http.StatusBadRequest)
		return
	}
	blockIds := r.PostForm["BlockIds"]
	err = writeFile(path, blockIds, app)
	httputil.WriteResponse(w, nil, wrr.Wrap(err))
}

func mergeFileHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("mergeFileHandler")
	out := r.PostFormValue("Out")
	if !fsutil.ValidateFilePath(out) {
		http.Error(w, wrr.SWrap("Invalid file path "+out), http.StatusBadRequest)
		return
	}
	in := r.PostForm["In"]
	for _, path := range in {
		if !fsutil.ValidateFilePath(path) {
			http.Error(w, wrr.SWrap("Invalid file path "+path), http.StatusBadRequest)
			return
		}
	}
	err := mergeFile(in, out)
	httputil.WriteResponse(w, nil, wrr.Wrap(err))
}
