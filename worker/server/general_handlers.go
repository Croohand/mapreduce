package server

import (
	"net/http"

	"github.com/Croohand/mapreduce/common/httputil"
)

func isAliveHandler(w http.ResponseWriter, r *http.Request) {
	if r.PostFormValue("Switch") == "true" {
		Config.CoordinatorAddr = r.PostFormValue("CoordinatorAddr")
	}
	httputil.WriteResponse(w, isAlive(), nil)
}
