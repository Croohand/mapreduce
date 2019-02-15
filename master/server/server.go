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

	http.HandleFunc("/IsAlive", isAlive)
	http.HandleFunc("/GetAvailableSlaves", getAvailableSlaves)
	http.HandleFunc("/GetMrConfig", getMRConfig)
	http.HandleFunc("/Transaction/IsAlive", isAliveTransaction)
	http.HandleFunc("/Transaction/Update", updateTransaction)
	http.HandleFunc("/Transaction/ValidateWrite", validateWriteTransaction)
	http.HandleFunc("/File/IsExists", isFileExists)

	log.Printf("starting master server with config %+v", Config)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil); err != nil {
		panic(err)
	}
}
