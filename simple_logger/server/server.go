package server

import (
	"fmt"
	"log"
	"net/http"
)

type LoggerConfig struct {
	Env        string
	Port       int
	Name       string
	OutputFile string
}

var Config LoggerConfig

func Run() {
	StartLogging()
	defer StopLogging()

	http.HandleFunc("/LogEntry", logEntryHandler)

	log.Printf("Starting logger server with config %+v", Config)
	addr := fmt.Sprintf(":%d", Config.Port)
	if Config.Env == "dev" {
		addr = "localhost" + addr
	}
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
