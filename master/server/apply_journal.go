package server

import (
	"bufio"
	"encoding/json"
	"log"
	"mime/multipart"
	"time"

	"github.com/Croohand/mapreduce/common/osutil"
	"github.com/Croohand/mapreduce/master/server/dbase"
)

func applyJournal(j multipart.File) {
	sc := bufio.NewScanner(j)
	for sc.Scan() {
		line := sc.Text()

		var e map[string]string
		err := json.Unmarshal([]byte(line), &e)
		if err != nil {
			log.Println(err)
			continue
		}

		if tsRaw, ok := e["ts"]; !ok {
			log.Println("No timestamp in journal entry")
			continue
		} else {
			var ts time.Time
			err := ts.UnmarshalText([]byte(tsRaw))
			if err != nil {
				log.Println(err)
				continue
			}
			if ts.Before(Config.LastJournalTs) {
				continue
			}
			Config.LastJournalTs = ts
		}

		err = dbase.ApplyEntry(e)

		if err != nil {
			log.Println(err)
			continue
		}
	}
	osutil.SaveConfig(Config)
}
