package commands

import "github.com/Croohand/mapreduce/common/httputil"

func startReadTransaction(paths []string) (string, *httputil.TransactionHandler) {
	return httputil.StartReadTransaction(mrConfig.GetHost(), paths, true)
}

func startWriteTransaction(paths []string) (string, *httputil.TransactionHandler) {
	return httputil.StartWriteTransaction(mrConfig.GetHost(), paths, true)
}
