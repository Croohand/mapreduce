package server

import (
	"net/http"
	"strconv"

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

func getAvailableReducersHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("getAvailableReducersHandler")
	num, err := strconv.Atoi(r.PostFormValue("Number"))
	if err != nil {
		http.Error(w, wrr.SWrap(err.Error()), http.StatusBadRequest)
		return
	}
	if num <= 0 {
		http.Error(w, wrr.SWrap("Invalid number of reducers "+strconv.Itoa(num)), http.StatusBadRequest)
		return
	}
	resp, err := getAvailableReducers(num)
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}

func getAvailableSchedulerHandler(w http.ResponseWriter, r *http.Request) {
	wrr := wrrors.New("getAvailableSchedulerHandler")
	resp, err := getAvailableScheduler()
	httputil.WriteResponse(w, resp, wrr.Wrap(err))
}

func getMrConfigHandler(w http.ResponseWriter, r *http.Request) {
	httputil.WriteResponse(w, getMrConfig(), nil)
}
