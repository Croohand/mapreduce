package server

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Croohand/mapreduce/common/fsutil"
)

func copyTemplate(txPath string) error {
	templatePath := filepath.Join("sources", "main")
	files, err := ioutil.ReadDir(templatePath)
	if err != nil {
		return err
	}
	for _, fileInfo := range files {
		filePath := filepath.Join(templatePath, fileInfo.Name())
		if strings.HasSuffix(fileInfo.Name(), ".go") {
			newPath := filepath.Join(txPath, "src", "main", fileInfo.Name())
			err := os.Link(filePath, newPath)
			if err != nil && os.IsNotExist(err) {
				return err
			}
		}
	}

	return nil
}

func buildSource(txId string) error {
	txPath, err := filepath.Abs(fsutil.GetTxDir(txId))
	if err != nil {
		return err
	}
	err = copyTemplate(txPath)
	if err != nil {
		return err
	}
	buildScriptPath := filepath.Join("sources", "build.sh")
	buildCmd := exec.Command("bash", buildScriptPath, txPath)
	var stderr bytes.Buffer
	buildCmd.Stderr = &stderr
	err = buildCmd.Run()
	if err != nil {
		return errors.New(err.Error() + ": " + stderr.String())
	}
	return nil
}
