package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/master/server/dbase"
)

type MasterConfig struct {
	Env            string
	Port           int
	Name           string
	SlaveAddrs     []string
	MasterAddrs    []string
	SchedulerAddrs []string
	LoggerAddr     string
	LastJournalTs  time.Time
}

var Config MasterConfig
var state = "passive"
var httpClient = httputil.NewClient("")

func RunServices() {
	log.Println("Starting master global processes")

	go monitorSlaves()
	go monitorTransactions()
	go monitorFiles()
	dbase.StartJournal(Config.MasterAddrs)
}

func Run() {
	log.Println("Opening bolt database")
	if Config.Env == "dev" {
		dbase.Open(Config.Name)
	} else {
		dbase.Open("")
	}
	defer dbase.Close()
	defer dbase.StopJournal()

	routes()

	go monitorMasters()

	log.Printf("Starting master server with config %+v", Config)
	addr := fmt.Sprintf(":%d", Config.Port)
	mux := http.Handler(http.DefaultServeMux)
	if Config.Env == "dev" {
		addr = "localhost" + addr
		mux = httputil.DefaultMuxWithLogging{Config.Name, Config.LoggerAddr}
		httpClient = httputil.NewClient(Config.Name)
	}
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
