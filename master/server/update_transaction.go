package server

import (
	"errors"
	"time"

	"github.com/Croohand/mapreduce/common/responses"
)

func updateTransaction(id, path string) (r *responses.Answer, err error) {
	r = &responses.Answer{}
	v, ok := transactions.Get(path)
	if ok {
		tx, ok := v.(Transaction)
		if !ok {
			transactions.Remove(path)
		} else if time.Since(tx.LastUpdate).Seconds() > transactionWaitTime {
			transactions.Remove(path)
			if tx.Id == id {
				err = errors.New("transaction is too old")
				return
			}
		} else if tx.Id != id {
			return
		}
	}
	transactions.Set(path, Transaction{Id: id, LastUpdate: time.Now()})
	r.Success = true
	return
}
