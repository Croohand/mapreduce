package server

import (
	"github.com/Croohand/mapreduce/common/fsutil"
)

func removeTransaction(txId string) error {
	return fsutil.RemoveTxDir(txId)
}
