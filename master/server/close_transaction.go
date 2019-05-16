package server

import (
	"errors"
	"fmt"
	"log"

	"github.com/Croohand/mapreduce/master/server/dbase"
)

func closeTransaction(txId string) error {
	has, err := dbase.Has(dbase.Txs, txId)
	if err != nil {
		return err
	}
	if !has {
		return errors.New(fmt.Sprintf("No transaction to close with id %s", txId))
	}
	tx, err := getTx(txId)
	if err != nil {
		return err
	}
	for _, path := range tx.Paths {
		txs, err := getPathTxs(path)
		if err != nil {
			log.Println(err)
			continue
		}
		if txs != nil {
			delete(txs, txId)
		}
		err = dbase.SetObject(dbase.PathTxs, path, txs)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	err = dbase.Del(dbase.Txs, txId)
	if err != nil {
		return err
	}
	return nil
}
