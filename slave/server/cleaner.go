package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/Croohand/mapreduce/common/httputil"
)

func logError(err error) {
	log.Println(Config.Name + " cleaner error: " + err.Error())
}

func cleaner() {
	for {
		log.Println(Config.Name + " cleaner starting new iteration")
		files, err := ioutil.ReadDir("transactions/")
		if err != nil {
			logError(err)
			continue
		}
		for _, f := range files {
			if f.IsDir() {
				id := f.Name()
				metaFile, err := os.Open(filepath.Join("transactions", id, "meta"))
				if err != nil {
					logError(err)
					continue
				}
				defer metaFile.Close()
				bytes, err := ioutil.ReadAll(metaFile)
				if err != nil {
					logError(err)
					continue
				}
				var meta struct{ Path string }
				if err := json.Unmarshal(bytes, &meta); err != nil {
					logError(err)
					continue
				}
				resp, err := http.PostForm(Config.MasterAddr+"/Transaction/IsAlive", url.Values{"Path": []string{meta.Path}, "TransactionId": []string{id}})
				if err != nil {
					logError(err)
					continue
				}
				var alive struct{ Alive bool }
				if err := httputil.GetJson(resp, &alive); err != nil {
					logError(err)
					continue
				}
				if !alive.Alive {
					if err := removeTransactionInner(id); err != nil {
						logError(err)
					}
				}
			}
		}
		time.Sleep(time.Minute*5 + time.Second*time.Duration(rand.Intn(300)))
	}
}
