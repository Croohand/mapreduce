package dbase

import (
	"os"
)

const (
	JournalPath = "journal.log"
)

var entries chan []byte
var exit chan bool
var masterAddrs []string
var started = false

func StartJournal(mAddrs []string) {
	started = true
	masterAddrs = mAddrs
	entries = make(chan []byte, 10000)
	exit = make(chan bool)
	go run()
}

func StopJournal() {
	exit <- true
}

func openJournal() *os.File {
	f, err := os.OpenFile(JournalPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return f
}

func run() {
	cur := 0
	f := openJournal()
	for {
		var e []byte

		select {
		case e = <-entries:
		case <-exit:
			f.Close()
			return
		}

		e = append(e, byte('\n'))
		_, err := f.Write(e)
		if err != nil {
			panic(err)
		}
		cur += 1
		if cur >= 10 {
			err := f.Close()
			if err != nil {
				panic(err)
			}
			cur = 0
			dumpJournal()
			f = openJournal()
		}
	}
}

func dumpJournal() {
}

func log(info []byte) {
	if started {
		entries <- info
	}
}
