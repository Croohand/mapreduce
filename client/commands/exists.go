package commands

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/Croohand/mapreduce/common/fsutil"
	"github.com/Croohand/mapreduce/common/httputil"
	"github.com/Croohand/mapreduce/common/responses"
)

func existsInner(path string) bool {
	if !fsutil.ValidateFilePath(path) {
		log.Panic("Invalid file path " + path)
	}
	resp, err := http.PostForm(mrConfig.Host+"/File/IsExists", url.Values{"Path": {path}})
	if err != nil {
		log.Panic(err)
	}
	var fStatus responses.FileStatus
	if err := httputil.GetJson(resp, &fStatus); err != nil {
		log.Panic(err)
	}
	return fStatus.Exists
}

func Exists(path string) {
	if existsInner(path) {
		fmt.Printf("File path %s exists\n", path)
	} else {
		fmt.Printf("File path %s doesn't exist\n", path)
	}
}
