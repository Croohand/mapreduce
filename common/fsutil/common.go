package fsutil

import (
	"os"
	"path/filepath"
)

func CreateTxDir(txId string) error {
	txPath := filepath.Join("transactions", txId)
	_, err := os.Stat(txPath)
	if os.IsNotExist(err) {
		if err := os.Mkdir(txPath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func RemoveTxDir(txId string) error {
	path := filepath.Join("transactions", txId)
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

func GetBlockPath(id, txId string) string {
	if txId == "" {
		return filepath.Join("blocks", id)
	}
	return filepath.Join("transactions", txId, id)
}
