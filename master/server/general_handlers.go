package server

import (
	"net/http"

	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/wrrors"
)

func isAliveHandler(w http.ResponseWriter, r *http.Request) {
	httputil.WriteResponse(w, isAlive(), nil)
}

func getAvailableSlavesHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("getAvailableSlavesHandler")
	resp, err := getAvailableSlaves(getMrConfig().ReplicationFactor)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}

func getMrConfigHandler(w http.ResponseWriter, r *http.Request) {
	httputil.WriteResponse(w, getMrConfig(), nil)
}
