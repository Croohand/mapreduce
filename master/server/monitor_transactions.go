package server

import (
	"log"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/timeutil"
	"github.com/Croohand/mapreduce/common/wrrors"
	"github.com/Croohand/mapreduce/master/server/dbase"
)

func monitorTransactions() {
	wrr := wrrors.New("monitorTransactions")
	for {
		timeutil.HugeSleep()
		log.Println(Config.Name + " monitorTransactions starting new iteration")
		txIds, err := dbase.GetKeys(dbase.Txs, "")
		if err != nil {
			log.Println(wrr.Wrap(err))
			continue
		}
		for _, txId := range txIds {
			var tx fsutil.Transaction
			err := dbase.GetObject(dbase.Txs, txId, &tx)
			if err != nil {
				log.Println(wrr.Wrap(err))
				continue
			}
			if !tx.IsAlive() {
				err = closeTransaction(txId)
				if err != nil {
					log.Println(wrr.Wrap(err))
				}
			}
		}
	}
}
