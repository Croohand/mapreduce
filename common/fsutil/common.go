package fsutil

import (
	"os"
	"path/filepath"
)

func CreateTxDir(txId string) error {
	txPath := GetTxDir(txId)
	_, err := os.Stat(txPath)
	if os.IsNotExist(err) {
		if err := os.Mkdir(txPath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func CreateShuffleDir(txId string) error {
	path := filepath.Join(GetTxDir(txId), "shuffle")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func RemoveTxDir(txId string) error {
	path := GetTxDir(txId)
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

func GetTxDir(txId string) string {
	return filepath.Join("transactions", txId)
}

func GetBlockPath(id, txId string) string {
	if txId == "" {
		return filepath.Join("blocks", id)
	}
	return filepath.Join("transactions", txId, id)
}

func CreateSourcesDir(txId string) error {
	srcsPath := filepath.Join("transactions", txId, "src", "mruserlib")
	return os.MkdirAll(srcsPath, os.ModePerm)
}

func CreateMainDir(txId string) error {
	mainPath := filepath.Join("transactions", txId, "src", "main")
	return os.MkdirAll(mainPath, os.ModePerm)
}

func GetSourcesDir(txId string) string {
	return filepath.Join("transactions", txId, "src", "mruserlib")
}

func GetSourcePath(name, txId string) string {
	return filepath.Join(GetSourcesDir(txId), name)
}
