package server

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

func readBlock(id string, w io.Writer) error {
	path := filepath.Join("files", id)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return errors.New("block with id " + id + " doesn't exist")
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}
	return nil
}
