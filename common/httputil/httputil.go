package httputil

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func WriteJson(w http.ResponseWriter, obj interface{}) {
	ans, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(ans); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetJson(r *http.Response, res interface{}) error {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return errors.New("GetJson: status " + r.Status + ": " + string(bytes))
	}
	return json.Unmarshal(bytes, res)
}
