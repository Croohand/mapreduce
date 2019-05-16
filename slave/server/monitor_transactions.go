package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
	"github.com/Croohand/mapreduce/common/timeutil"
	"github.com/Croohand/mapreduce/common/wrrors"
)

func monitorTransactions() {
	wrr := wrrors.New("monitorTransactions")
	for {
		timeutil.HugeSleep()
		log.Println(Config.Name + " monitorTransactions starting new iteration")
		files, err := ioutil.ReadDir("transactions/")
		if err != nil {
			log.Println(wrr.Wrap(err))
			continue
		}
		for _, f := range files {
			if f.IsDir() {
				txId := f.Name()
				resp, err := http.PostForm(Config.MasterAddr+"/Transaction/IsAlive", url.Values{"TransactionId": {txId}})
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
					if err := removeTransaction(txId); err != nil {
						log.Println(wrr.Wrap(err))
					}
				}
			}
		}
	}
}
