package server

import "os"

var entries chan string
var exit chan bool

func StartLogging() {
	entries = make(chan string, 10000)
	exit = make(chan bool)
	go run()
}

func StopLogging() {
	exit <- true
}

func logEntry(e string) {
	entries <- e
}

func run() {
	f, err := os.OpenFile(Config.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for {
		select {
		case e := <-entries:
			f.Write([]byte(e + "\n"))
		case <-exit:
			return
		}
	}
}
