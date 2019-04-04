package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	cmap "github.com/orcaman/concurrent-map"
	bolt "go.etcd.io/bbolt"
)

type Transaction struct {
	Id         string
	Path       string
	LastUpdate time.Time
}

var transactions = cmap.New()
var filesDB *bolt.DB

type MasterConfig struct {
	Port       int
	Name       string
	SlaveAddrs []string
}

var Config MasterConfig

func Run() {
	log.Println("opening bolt database")
	var err error
	filesDB, err = bolt.Open("files.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer filesDB.Close()

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
