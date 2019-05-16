package server

import (
	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/master/server/dbase"
)

func mergeFile(in []string, out string) error {
	var resIds fsutil.BlockIds
	for _, path := range in {
		var curIds fsutil.BlockIds
		err := dbase.GetObject(dbase.Files, path, &curIds)
		if err != nil {
			return err
		}
		resIds = append(resIds, curIds...)
	}
	return dbase.SetObject(dbase.Files, out, resIds)
}
