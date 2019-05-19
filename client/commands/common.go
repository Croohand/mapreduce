package commands

import "github.com/Croohand/mapreduce/common/httputil"

func startReadTransaction(paths []string) (string, *httputil.TransactionHandler) {
	return httputil.StartReadTransaction(mrConfig.Host, paths, true)
}

func startWriteTransaction(paths []string) (string, *httputil.TransactionHandler) {
	return httputil.StartWriteTransaction(mrConfig.Host, paths, true)
}
