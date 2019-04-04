package server

import (
	"net/http"

	"github.com/Croohand/mapreduce/common/httputil"
)

func isAliveHandler(w http.ResponseWriter, r *http.Request) {
	httputil.WriteResponse(w, isAlive(), nil)
}
