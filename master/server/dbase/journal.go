package dbase

import (
	"bytes"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/Croohand/mapreduce/common/httputil"
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
	if started {
		exit <- true
	}
}

func ApplyEntry(e map[string]string) error {
	t, ok := e["type"]
	if !ok {
		return errors.New("No type in journal entry")
	}
	var err error
	switch t {
	case "set":
		err = Set(e["bucket"], e["key"], []byte(e["value"]))
	case "del":
		err = Del(e["bucket"], e["key"])
	default:
		return errors.New("Unknown type " + t + " in journal entry")
	}
	if err != nil {
		return err
	}
	return nil
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
	var orig bytes.Buffer
	f, err := os.Open(JournalPath)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = io.Copy(&orig, f)
	if err != nil {
		log.Println(err)
		return
	}
	f.Close()

	success := 0

	for _, addr := range masterAddrs[1:] {
		cur := orig

		var b bytes.Buffer
		w := multipart.NewWriter(&b)

		fw, err := w.CreateFormFile("Journal", "File")
		if err != nil {
			continue
		}
		if _, err = io.Copy(fw, &cur); err != nil {
			continue
		}

		w.Close()

		req, err := http.NewRequest("POST", addr+"/Journal/Apply", &b)
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.Close = true
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		if err := httputil.GetError(resp); err != nil {
			continue
		}
		success += 1
	}

	if success == len(masterAddrs)-1 {
		f, err := os.Create(JournalPath)
		if err != nil {
			f.Close()
		}
	}
}

func logJournal(info []byte) {
	if started {
		entries <- info
	}
}
