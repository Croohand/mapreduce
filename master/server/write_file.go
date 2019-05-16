package server

import (
	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/master/server/dbase"
)

func writeFile(path string, blockIds fsutil.BlockIds, app bool) error {
	if app {
		has, err := dbase.Has(dbase.Files, path)
		if err != nil {
			return err
		}
		if has {
			var prevBlockIds fsutil.BlockIds
			err := dbase.GetObject(dbase.Files, path, &prevBlockIds)
			if err != nil {
				return err
			}
			blockIds = append(prevBlockIds, blockIds...)
		}
	}
	return dbase.SetObject(dbase.Files, path, blockIds)
}
