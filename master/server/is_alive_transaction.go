package server

import (
	"time"

	"github.com/Croohand/mapreduce/common/responses"
)

const transactionWaitTime = 60

func isAliveTransaction(id, path string) (r *responses.TransactionStatus) {
	r = &responses.TransactionStatus{}
	v, ok := transactions.Get(path)
	r.Alive = false
	if ok {
		tx, ok := v.(Transaction)
		if !ok || time.Since(tx.LastUpdate).Seconds() > transactionWaitTime {
			transactions.Remove(path)
			return
		}
		r.Alive = (tx.Id == id)
	}
	return
}
