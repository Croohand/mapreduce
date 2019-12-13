package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/Croohand/mapreduce/common/httputil"
)

type SlaveConfig struct {
	Env        string
	Port       int
	Name       string
	MasterAddr string
	LoggerAddr string
	Scheduler  bool
}

var Config SlaveConfig
var httpClient = httputil.NewClient("")

func logStatus(status string) {
	if Config.Env != "dev" {
		return
	}
	e := fmt.Sprintf("%v %v", Config.Name, status)
	httpClient.PostForm(Config.LoggerAddr+"/LogEntry", url.Values{"Entry": {e}})
}

func cleanup() {
	logStatus("down")
}

func Run() {
	if !Config.Scheduler {
		log.Println("Creating blocks directory for slave " + Config.Name)
		if err := os.Mkdir("blocks", os.ModePerm); err != nil && !os.IsExist(err) {
			panic(err)
		}
		log.Println("Checking sources directory existence for slave " + Config.Name)
		_, err := os.Stat("sources")
		if err != nil && !os.IsExist(err) {
			panic(err)
		}
	}

	log.Println("Creating transactions directory for slave " + Config.Name)
	if err := os.Mkdir("transactions", os.ModePerm); err != nil && !os.IsExist(err) {
		panic(err)
	}

	if !Config.Scheduler {
		http.HandleFunc("/Block/Check", checkBlockHandler)
		http.HandleFunc("/Block/Write", writeBlockHandler)
		http.HandleFunc("/Block/Copy", copyBlockHandler)
		http.HandleFunc("/Block/Remove", removeBlockHandler)
		http.HandleFunc("/Block/Read", readBlockHandler)
		http.HandleFunc("/Block/Validate", validateBlockHandler)
		http.HandleFunc("/Transaction/Remove", removeTransactionHandler)
		http.HandleFunc("/Source/Build", buildSourceHandler)
		http.HandleFunc("/Operation/Map", mapOperationHandler)
		http.HandleFunc("/Operation/Reduce", reduceOperationHandler)
		http.HandleFunc("/Operation/SendResults", sendResultsOperationHandler)
	} else {
		http.HandleFunc("/Operation/PrepareMapReduce", prepareMapReduceOperationHandler)
		http.HandleFunc("/Operation/StartMapReduce", startMapReduceOperationHandler)
		http.HandleFunc("/Operation/GetStatus", getStatusOperationHandler)
	}

	http.HandleFunc("/IsAlive", isAliveHandler)
	http.HandleFunc("/Source/Write", writeSourceHandler)

	log.Printf("Starting global processes on slave server")
	go monitorTransactions()

	log.Printf("Starting slave server with config %+v", Config)
	addr := fmt.Sprintf(":%d", Config.Port)
	mux := http.Handler(http.DefaultServeMux)
	if Config.Env == "dev" {
		addr = "localhost" + addr
		mux = httputil.DefaultMuxWithLogging{Config.Name, Config.LoggerAddr}
		httpClient = httputil.NewClient(Config.Name)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	logStatus("up")
	defer cleanup()
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
