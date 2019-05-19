package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Croohand/mapreduce/master/server/dbase"
)

type MasterConfig struct {
	Port           int
	Name           string
	SlaveAddrs     []string
	SchedulerAddrs []string
}

var Config MasterConfig

func Run() {
	log.Println("Opening bolt database")
	dbase.Open()
	defer dbase.Close()

	http.HandleFunc("/IsAlive", isAliveHandler)
	http.HandleFunc("/GetAvailableSlaves", getAvailableSlavesHandler)
	http.HandleFunc("/GetAvailableScheduler", getAvailableSchedulerHandler)
	http.HandleFunc("/GetMrConfig", getMrConfigHandler)
	http.HandleFunc("/Transaction/IsAlive", isAliveTransactionHandler)
	http.HandleFunc("/Transaction/Update", updateTransactionHandler)
	http.HandleFunc("/Transaction/Start", startTransactionHandler)
	http.HandleFunc("/Transaction/Close", closeTransactionHandler)
	http.HandleFunc("/Transaction/ValidateBlocks", validateBlocksHandler)
	http.HandleFunc("/File/IsExists", isFileExistsHandler)
	http.HandleFunc("/File/Remove", removeFileHandler)
	http.HandleFunc("/File/Write", writeFileHandler)
	http.HandleFunc("/File/Read", readFileHandler)
	http.HandleFunc("/File/Merge", mergeFileHandler)
	http.HandleFunc("/File/List", listFilesHandler)

	log.Println("Starting master global processes")

	go monitorSlaves()
	go monitorTransactions()
	go monitorFiles()

	log.Printf("Starting master server with config %+v", Config)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil); err != nil {
		panic(err)
	}
}
