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
	"github.com/Croohand/mapreduce/common/responses"
	"github.com/Croohand/mapreduce/common/wrrors"
)

func cleaner() {
	wrr := wrrors.New("cleaner")
	for {
		log.Println(Config.Name + " cleaner starting new iteration")
		files, err := ioutil.ReadDir("transactions/")
		if err != nil {
			log.Println(wrr.Wrap(err))
			continue
		}
		for _, f := range files {
			if f.IsDir() {
				id := f.Name()
				metaFile, err := os.Open(filepath.Join("transactions", id, "meta"))
				if err != nil {
					log.Println(wrr.Wrap(err))
					continue
				}
				defer metaFile.Close()
				bytes, err := ioutil.ReadAll(metaFile)
				if err != nil {
					log.Println(wrr.Wrap(err))
					continue
				}
				var meta struct{ Path string }
				if err := json.Unmarshal(bytes, &meta); err != nil {
					log.Println(wrr.Wrap(err))
					continue
				}
				resp, err := http.PostForm(Config.MasterAddr+"/Transaction/IsAlive", url.Values{"Path": {meta.Path}, "TransactionId": {id}})
				if err != nil {
					log.Println(wrr.Wrap(err))
					continue
				}
				var txStatus responses.TransactionStatus
				if err := httputil.GetJson(resp, &txStatus); err != nil {
					log.Println(wrr.Wrap(err))
					continue
				}
				if !txStatus.Alive {
					if _, err := removeTransaction(id); err != nil {
						log.Println(wrr.Wrap(err))
					}
				}
			}
		}
		time.Sleep(time.Minute*5 + time.Second*time.Duration(rand.Intn(300)))
	}
}
