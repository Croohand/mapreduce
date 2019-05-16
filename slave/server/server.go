package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type SlaveConfig struct {
	Port       int
	Name       string
	MasterAddr string
}

var Config SlaveConfig

func Run() {
	log.Println("Creating blocks directory for slave " + Config.Name)
	if err := os.Mkdir("blocks", os.ModePerm); err != nil && !os.IsExist(err) {
		panic(err)
	}
	log.Println("Creating transactions directory for slave " + Config.Name)
	if err := os.Mkdir("transactions", os.ModePerm); err != nil && !os.IsExist(err) {
		panic(err)
	}

	http.HandleFunc("/IsAlive", isAliveHandler)
	http.HandleFunc("/Block/Check", checkBlockHandler)
	http.HandleFunc("/Block/Write", writeBlockHandler)
	http.HandleFunc("/Block/Copy", copyBlockHandler)
	http.HandleFunc("/Block/Remove", removeBlockHandler)
	http.HandleFunc("/Block/Read", readBlockHandler)
	http.HandleFunc("/Block/Validate", validateBlockHandler)
	http.HandleFunc("/Transaction/Remove", removeTransactionHandler)

	log.Printf("Starting global processes on slave server")
	go monitorTransactions()

	log.Printf("Starting slave server with config %+v", Config)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil); err != nil {
		panic(err)
	}
}
