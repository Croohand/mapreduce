package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Croohand/mapreduce/master/server/dbase"
	cmap "github.com/orcaman/concurrent-map"
)

type Transaction struct {
	Id         string
	Path       string
	LastUpdate time.Time
}

var transactions = cmap.New()

type MasterConfig struct {
	Port       int
	Name       string
	SlaveAddrs []string
}

var Config MasterConfig

func Run() {
	log.Println("opening bolt database")
	dbase.Open()
	defer dbase.Close()

	http.HandleFunc("/IsAlive", isAliveHandler)
	http.HandleFunc("/GetAvailableSlaves", getAvailableSlavesHandler)
	http.HandleFunc("/GetMrConfig", getMrConfigHandler)
	http.HandleFunc("/Transaction/IsAlive", isAliveTransactionHandler)
	http.HandleFunc("/Transaction/Update", updateTransactionHandler)
	http.HandleFunc("/Transaction/ValidateWrite", validateWriteTransactionHandler)
	http.HandleFunc("/File/IsExists", isFileExistsHandler)

	log.Printf("starting master server with config %+v", Config)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil); err != nil {
		panic(err)
	}
}
