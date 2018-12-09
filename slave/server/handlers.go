package server

import (
	"net/http"

	"github.com/Croohand/mapreduce/common/httputil"
)

func isAlive(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJson(w, struct {
		Alive bool
		Type  string
	}{true, "slave"})
}
