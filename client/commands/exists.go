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

func Exists(path string) bool {
	if !fsutil.ValidateFilePath(path) {
		log.Fatal("invalid file path " + path)
	}
	resp, err := http.PostForm(mrConfig.Host+"/File/IsExists", url.Values{"Path": {path}})
	if err != nil {
		log.Fatal(err)
	}
	var fStatus responses.FileStatus
	if err := httputil.GetJson(resp, &fStatus); err != nil {
		log.Fatal(err)
	}
	if fStatus.Exists {
		fmt.Printf("File %s exists\n", path)
	} else {
		fmt.Printf("File %s doesn't exist\n", path)
	}
	return fStatus.Exists
}
