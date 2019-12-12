package server

import "net/http"

type handlerFunc func(http.ResponseWriter, *http.Request)

var (
	routesPassive map[string]handlerFunc = map[string]handlerFunc{
		"IsAlive":               isAliveHandler,
		"GetAvailableSlaves":    getAvailableSlavesHandler,
		"GetAvailableScheduler": getAvailableSchedulerHandler,
		"GetAvailableReducers":  getAvailableReducersHandler,
		"GetMrConfig":           getMrConfigHandler,
		"Transaction/IsAlive":   isAliveTransactionHandler,
		"File/IsExists":         isFileExistsHandler,
		"File/Read":             readFileHandler,
		"File/List":             listFilesHandler,
		"Journal/Apply":         applyJournalHandler,
	}

	routesActive map[string]handlerFunc = map[string]handlerFunc{
		"Transaction/Update":         updateTransactionHandler,
		"Transaction/Start":          startTransactionHandler,
		"Transaction/Close":          closeTransactionHandler,
		"Transaction/ValidateBlocks": validateBlocksHandler,
		"File/Remove":                removeFileHandler,
		"File/Write":                 writeFileHandler,
		"File/Merge":                 mergeFileHandler,
	}
)

func checkState(f handlerFunc) handlerFunc {
	res := func(w http.ResponseWriter, r *http.Request) {
		if state == "passive" {
			return
		}
		f(w, r)
	}
	return res
}

func routes() {
	for path, handler := range routesPassive {
		http.HandleFunc("/"+path, handler)
	}

	for path, handler := range routesActive {
		http.HandleFunc("/"+path, checkState(handler))
	}
}
