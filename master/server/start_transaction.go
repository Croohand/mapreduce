package server

import (
	"errors"
	"fmt"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/responses"
	"github.com/Croohand/mapreduce/master/server/dbase"
)

func getPathTxs(path string) (fsutil.TxIds, error) {
	has, err := dbase.Has(dbase.PathTxs, path)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	var txs fsutil.TxIds
	return txs, dbase.GetObject(dbase.PathTxs, path, &txs)
}

func getTx(txId string) (*fsutil.Transaction, error) {
	var tx fsutil.Transaction
	err := dbase.GetObject(dbase.Txs, txId, &tx)
	return &tx, err
}

func startTransaction(paths []string, txType fsutil.TransactionType) (*responses.StartedTransaction, error) {
	if len(paths) == 0 {
		return nil, errors.New("Can't start transaction on no paths")
	}
	for _, path := range paths {
		txs, err := getPathTxs(path)
		if err != nil {
			return nil, err
		}
		if txs != nil {
			for txId := range txs {
				tx, err := getTx(txId)
				if err != nil {
					return nil, err
				}
				if !tx.IsAlive() {
					closeTransaction(txId)
				} else if tx.TxType == fsutil.TxTypeWrite || txType == fsutil.TxTypeWrite {
					return nil, errors.New(fmt.Sprintf("Couldn't take transaction for %s with type %v because of concurrent transaction with id %s and type %v", path, txType, txId, tx.TxType))
				}
			}
		}
	}
	tx := fsutil.NewTransaction(paths, txType)
	txId := fsutil.GenerateTransactionId()
	err := dbase.SetObject(dbase.Txs, txId, tx)
	if err != nil {
		return nil, err
	}
	for _, path := range paths {
		txs, err := getPathTxs(path)
		if err != nil {
			return nil, err
		}
		if txs == nil {
			txs = make(fsutil.TxIds)
		}
		txs[txId] = struct{}{}
		err = dbase.SetObject(dbase.PathTxs, path, txs)
		if err != nil {
			return nil, err
		}
	}
	return &responses.StartedTransaction{txId}, nil
}
