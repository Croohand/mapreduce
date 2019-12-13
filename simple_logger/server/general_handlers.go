package server

import (
	"net/http"
	"time"
)

func logEntryHandler(w http.ResponseWriter, r *http.Request) {
	e := r.PostFormValue("Entry")
	if e == "" {
		http.Error(w, "Error empty entry", http.StatusBadRequest)
		return
	}
	logEntry(time.Now().Format("15:04:05.9999999") + " " + e)
}
