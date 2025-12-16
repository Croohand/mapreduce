package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/coordinator/server/dbase"
)

type CoordinatorConfig struct {
	Env            string
	Port           int
	Name           string
	WorkerAddrs     []string
	CoordinatorAddrs    []string
	SchedulerAddrs []string
	LoggerAddr     string
	LastJournalTs  time.Time
}

var Config CoordinatorConfig
var state = "passive"
var httpClient = httputil.NewClient("")

func RunServices() {
	log.Println("Starting coordinator global processes")

	go monitorWorkers()
	go monitorTransactions()
	go monitorFiles()
	dbase.StartJournal(Config.CoordinatorAddrs)
}

func logStatus(status string) {
	if Config.Env != "dev" {
		return
	}
	e := fmt.Sprintf("%v %v", Config.Name, status)
	httpClient.PostForm(Config.LoggerAddr+"/LogEntry", url.Values{"Entry": {e}})
}

func cleanup() {
	dbase.Close()
	dbase.StopJournal()
	logStatus("down")
}

func Run() {
	log.Println("Opening bolt database")
	if Config.Env == "dev" {
		dbase.Open(Config.Name)
	} else {
		dbase.Open("")
	}

	routes()

	go monitorCoordinators()

	log.Printf("Starting coordinator server with config %+v", Config)
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
