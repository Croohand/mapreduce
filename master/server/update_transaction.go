package server

import (
	"errors"
	"fmt"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/master/server/dbase"
)

func updateTransaction(txId string) error {
	var tx fsutil.Transaction
	err := dbase.GetObject(dbase.Txs, txId, &tx)
	if err != nil {
		return err
	}
	if !tx.IsAlive() {
		return errors.New(fmt.Sprintf("Transaction with id %s is not alive", txId))
	}
	tx.Update()
	err = dbase.SetObject(dbase.Txs, txId, tx)
	return err
}
