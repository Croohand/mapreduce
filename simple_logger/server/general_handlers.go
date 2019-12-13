package server

import "net/http"

func logEntryHandler(w http.ResponseWriter, r *http.Request) {
	e := r.PostFormValue("Entry")
	if e == "" {
		http.Error(w, "Error empty entry", http.StatusBadRequest)
		return
	}
	logEntry(e)
}
