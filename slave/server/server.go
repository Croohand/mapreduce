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
	log.Println("creating files directory for slave " + Config.Name)
	if err := os.Mkdir("files", os.ModePerm); err != nil && !os.IsExist(err) {
		panic(err)
	}
	log.Println("creating transactions directory for slave " + Config.Name)
	if err := os.Mkdir("transactions", os.ModePerm); err != nil && !os.IsExist(err) {
		panic(err)
	}

	http.HandleFunc("/IsAlive", isAliveHandler)
	http.HandleFunc("/Block/IsExists", checkBlockHandler)
	http.HandleFunc("/Block/Write", writeBlockHandler)
	http.HandleFunc("/Block/Remove", removeBlockHandler)
	http.HandleFunc("/Block/Read", readBlockHandler)
	http.HandleFunc("/Block/Validate", validateBlockHandler)
	http.HandleFunc("/Transaction/Remove", removeTransactionHandler)

	log.Printf("starting cleaner on slave server")
	go cleaner()

	log.Printf("starting slave server with config %+v", Config)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil); err != nil {
		panic(err)
	}
}
