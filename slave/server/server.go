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

	http.HandleFunc("/IsAlive", isAlive)
	http.HandleFunc("/Block/IsExists", checkBlock)
	http.HandleFunc("/Block/Write", writeBlock)
	http.HandleFunc("/Block/Remove", removeBlock)
	http.HandleFunc("/Block/Read", readBlock)
	http.HandleFunc("/Block/Validate", validateBlock)
	http.HandleFunc("/Transaction/Remove", removeTransaction)

	log.Printf("starting cleaner on slave server")
	go cleaner()

	log.Printf("starting slave server with config %+v", Config)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil); err != nil {
		panic(err)
	}
}
