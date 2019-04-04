package server

import (
	"net/http"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/wrrors"
)

func isFileExistsHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("isFileExistsHandler")
	path := r.PostFormValue("Path")
	if !fsutil.ValidateFilePath(path) {
		http.Error(w, wrr.SWrap("invalid file path"), http.StatusBadRequest)
		return
	}
	resp, err := isFileExists(path)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}
