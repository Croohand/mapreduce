package httputil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Croohand/mapreduce/common/wrrors"
)

func writeJson(w http.ResponseWriter, obj interface{}) {
	ans, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(ans); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WriteResponse(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJson(w, resp)
}

func GetJson(r *http.Response, res interface{}) error {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return wrrors.New("GetJson").WrapS("status " + r.Status + ": " + string(bytes))
	}
	return json.Unmarshal(bytes, res)
}
