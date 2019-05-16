package server

import (
	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/responses"
	"github.com/Croohand/mapreduce/master/server/dbase"
)

func isAliveTransaction(txId string) (*responses.TransactionStatus, error) {
	has, err := dbase.Has(dbase.Txs, txId)
	if err != nil {
		return nil, err
	}
	if !has {
		return &responses.TransactionStatus{false}, nil
	}
	var tx fsutil.Transaction
	err = dbase.GetObject(dbase.Txs, txId, &tx)
	if err != nil {
		return nil, err
	}
	return &responses.TransactionStatus{tx.IsAlive()}, nil
}
