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

var mrConfig = struct {
	BlockSize            int
	ReplicationFactor    int
	MinReplicationFactor int
}{BlockSize: 1 << 20, ReplicationFactor: 3, MinReplicationFactor: 2}

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
	http.HandleFunc("/GetMRConfig", getMRConfig)
	http.HandleFunc("/Transaction/IsAlive", isAliveTransaction)
	http.HandleFunc("/Transaction/Update", updateTransaction)
	http.HandleFunc("/Transaction/ValidateWrite", validateWriteTransaction)

	log.Printf("starting master server with config %+v", Config)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil); err != nil {
		panic(err)
	}
}
