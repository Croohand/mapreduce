package server

import (
	"errors"

	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
)

func prepareMapReduceOperation(in []string, out string) (resp *responses.PreparedOperation, err error) {
	defer func() {
		if e := recover(); e != nil {
			switch er := e.(type) {
			case error:
				err = er
			case string:
				err = errors.New(er)
			}
		}
	}()
	readId, _ := httputil.StartReadTransaction(Config.MasterAddr, in, false)
	writeId, _ := httputil.StartWriteTransaction(Config.MasterAddr, []string{out}, false)
	resp = &responses.PreparedOperation{readId, writeId}
	return
}
